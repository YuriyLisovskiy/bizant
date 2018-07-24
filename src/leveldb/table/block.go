// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package table

import (
	"io"
	"sort"
	"bytes"
	"encoding/binary"
	"github.com/YuriyLisovskiy/blockchain-go/src/leveldb"
)

const maxInt = int(^uint(0) >> 1)

type blockHandle struct {
	Offset, Size uint64
}

func (p *blockHandle) DecodeFrom(b []byte) (int, error) {
	var n, m int
	p.Offset, n = binary.Uvarint(b)
	if n > 0 {
		p.Size, m = binary.Uvarint(b[n:])
	}
	if n <= 0 || m <= 0 {
		return 0, leveldb.NewCorruptionError("bad block handle")
	}
	return n + m, nil
}

func (p *blockHandle) Encode() []byte {
	b := make([]byte, binary.MaxVarintLen64*2)
	n := binary.PutUvarint(b, p.Offset)
	return b[:n+binary.PutUvarint(b[n:], p.Size)]
}

func (p *blockHandle) EncodeTo(b []byte) int {
	n := binary.PutUvarint(b, p.Offset)
	return n + binary.PutUvarint(b[n:], p.Size)
}

func readFullAt(r io.ReaderAt, buf []byte, off int64) (n int, err error) {
	for n < len(buf) && err == nil {
		var nn int
		nn, err = r.ReadAt(buf[n:], off+int64(n))
		n += nn
	}
	if err == io.EOF {
		if n >= len(buf) {
			err = nil
		} else if n > 0 {
			err = io.ErrUnexpectedEOF
		}
	}
	return n, err
}

func (p *blockHandle) ReadAll(r io.ReaderAt, checksum bool) (b []byte, err error) {
	raw := make([]byte, p.Size+5)
	_, err = readFullAt(r, raw, int64(p.Offset))
	if err != nil {
		return
	}
	crcb := raw[len(raw)-4:]
	raw = raw[:len(raw)-4]
	if checksum {
		var sum uint32
		bb := bytes.NewBuffer(crcb)
		err = binary.Read(bb, binary.LittleEndian, &sum)
		if err != nil {
			return
		}
		sum = leveldb.UnmaskCRC32(sum)
		crc := leveldb.NewCRC32C()
		crc.Write(raw)
		if crc.Sum32() != sum {
			err = leveldb.NewCorruptionError("block checksum mismatch")
			return
		}
	}
	compression := leveldb.Compression(raw[len(raw)-1])
	b = raw[:len(raw)-1]
	switch compression {
	case leveldb.SnappyCompression:
		compression = leveldb.NoCompression
	}
	return
}

type block struct {
	buf          []byte
	br, rr       *bytes.Reader
	restartLen   int
	restartStart int
}

func newBlock(buf []byte) (b *block, err error) {
	if len(buf) < 8 {
		err = leveldb.NewCorruptionError("block to short")
		return
	}
	br := bytes.NewReader(buf)
	var restartLen uint32
	br.Seek(int64(len(buf))-4, 0)
	err = binary.Read(br, binary.LittleEndian, &restartLen)
	if err != nil {
		return
	}
	restartStart := len(buf) - (1+int(restartLen))*4
	if restartStart > len(buf)-4 {
		err = leveldb.NewCorruptionError("bad restart offset in block")
		return
	}
	rr := bytes.NewReader(buf[restartStart:len(buf)-4])
	b = &block{buf: buf, br: br, rr: rr, restartLen: int(restartLen), restartStart: restartStart}
	return
}

func newBlockFromHandle(handle *blockHandle, r leveldb.Reader, verify bool) (b *block, err error) {
	var buf []byte
	buf, err = handle.ReadAll(r, verify)
	if err != nil {
		return nil, err
	}
	return newBlock(buf)
}

func (b *block) NewIterator(cmp leveldb.Comparator) leveldb.Iterator {
	if b.restartLen == 0 {
		return &leveldb.EmptyIterator{}
	}
	return &blockIterator{b: b, cmp: cmp}
}

func (b *block) getRestartOffset(idx int) (offset int, err error) {
	if idx >= b.restartLen {
		panic("out of range")
	}
	_, err = b.rr.Seek(int64(idx)*4, 0)
	if err != nil {
		return
	}
	var tmp uint32
	err = binary.Read(b.rr, binary.LittleEndian, &tmp)
	return int(tmp), err
}

func (b *block) getRestartRange(idx int) (r *restartRange, err error) {
	var start, end int
	start, err = b.getRestartOffset(idx)
	if err != nil {
		return
	}
	if start >= b.restartStart {
		goto corrupt
	}
	if idx+1 < b.restartLen {
		end, err = b.getRestartOffset(idx + 1)
		if err != nil {
			return
		}
		if end >= b.restartStart {
			goto corrupt
		}
	} else {
		end = b.restartStart
	}
	if start < end {
		r = &restartRange{raw: b.buf[start:end]}
		r.buf = bytes.NewBuffer(r.raw)
		return
	}
corrupt:
	return nil, leveldb.NewCorruptionError("bad restart range in block")
}

type keyVal struct {
	key, value []byte
}

type restartRange struct {
	raw    []byte
	buf    *bytes.Buffer
	kv     keyVal
	pos    int
	cached bool
	cache  []keyVal
}

func (r *restartRange) Next() (err error) {
	if r.cached && len(r.cache) > r.pos {
		r.kv = r.cache[r.pos]
		r.pos++
		return
	}
	if r.buf.Len() == 0 {
		return io.EOF
	}
	var nKey []byte
	var shared, nonShared, valueLen uint64
	shared, err = binary.ReadUvarint(r.buf)
	if err != nil || shared > uint64(len(r.kv.key)) {
		goto corrupt
	}
	nonShared, err = binary.ReadUvarint(r.buf)
	if err != nil {
		goto corrupt
	}
	valueLen, err = binary.ReadUvarint(r.buf)
	if err != nil {
		goto corrupt
	}
	if nonShared+valueLen > uint64(r.buf.Len()) {
		goto corrupt
	}
	if r.cached && r.pos > 0 {
		r.cache = append(r.cache, r.kv)
	}
	nKey = r.buf.Next(int(nonShared))
	if shared == 0 {
		r.kv.key = nKey
	} else {
		pKey, key := r.kv.key[:shared], make([]byte, shared+nonShared)
		copy(key, pKey)
		copy(key[shared:], nKey)
		r.kv.key = key
	}
	r.kv.value = r.buf.Next(int(valueLen))
	r.pos++
	return
corrupt:
	return leveldb.NewCorruptionError("bad entry in block")
}

func (r *restartRange) Prev() (err error) {
	if r.pos <= 1 {
		r.pos = 0
		return io.EOF
	}
	r.pos--
	if r.cached {
		if len(r.cache) == r.pos {
			r.cache = append(r.cache, r.kv)
		}
		r.kv = r.cache[r.pos-1]
		return err
	}
	r.cached = true
	pos := r.pos
	r.Reset()
	for r.pos < pos {
		err = r.Next()
		if err != nil {
			if err == io.EOF {
				panic("not reached")
			}
			return err
		}
	}
	return err
}

func (r *restartRange) Last() (err error) {
	if !r.cached {
		r.cached = true
		r.Reset()
	}
	for {
		err = r.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return
}

func (r *restartRange) Reset() {
	if r.pos > 0 {
		r.buf, r.pos = bytes.NewBuffer(r.raw), 0
	}
}

func (r *restartRange) Key() []byte {
	return r.kv.key
}

func (r *restartRange) Value() []byte {
	return r.kv.value
}

type blockIterator struct {
	b   *block
	cmp leveldb.Comparator
	err error
	ri  int
	rr  *restartRange
}

func (i *blockIterator) First() bool {
	i.ri, i.rr = 0, nil
	return i.Next()
}

func (i *blockIterator) Last() bool {
	i.ri, i.rr = i.b.restartLen, nil
	return i.Prev()
}

func (i *blockIterator) Seek(key []byte) (r bool) {
	if i.err != nil {
		return false
	}
	i.ri = sort.Search(i.b.restartLen, func(x int) bool {
		rr, err := i.b.getRestartRange(x)
		if err != nil {
			i.err = err
			panic(i)
		}
		err = rr.Next()
		if err != nil {
			i.err = err
			panic(i)
		}
		return i.cmp.Compare(rr.Key(), key) > 0
	})
	if i.ri > 0 {
		i.ri--
	}
	i.rr = nil
	for i.Next() {
		if i.cmp.Compare(i.rr.Key(), key) >= 0 {
			return true
		}
	}
	return false
}

func (i *blockIterator) Next() bool {
	if i.ri == i.b.restartLen || i.err != nil {
		return false
	}
	if i.rr == nil {
		i.setRestartRange()
		if i.err != nil {
			return false
		}
	}
	i.err = i.rr.Next()
	if i.err != nil {
		if i.err == io.EOF {
			i.err = nil
			i.ri++
			i.rr = nil
			return i.Next()
		}
		return false
	}
	return true
}

func (i *blockIterator) Prev() bool {
	if i.err != nil {
		return false
	}
	if i.rr == nil {
		if i.ri == 0 {
			return false
		}
		i.ri--
		i.setRestartRange()
		if i.err != nil {
			return false
		}
		i.err = i.rr.Last()
		if i.err != nil {
			return false
		}
		return true
	}
	i.err = i.rr.Prev()
	if i.err != nil {
		if i.err == io.EOF {
			i.err = nil
			i.rr = nil
			return i.Prev()
		}
		return false
	}

	return true
}

func (i *blockIterator) Key() []byte {
	if i.rr == nil {
		return nil
	}
	return i.rr.Key()
}

func (i *blockIterator) Value() []byte {
	if i.rr == nil {
		return nil
	}
	return i.rr.Value()
}

func (i *blockIterator) Error() error { return i.err }

func (i *blockIterator) setRestartRange() { i.rr, i.err = i.b.getRestartRange(i.ri) }
