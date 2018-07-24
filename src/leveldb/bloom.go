// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import (
	"io"
)

func BloomHash(key []byte) uint32 {
	return Hash(key, 0xbc9f1d34)
}

type BloomFilter struct {
	bitsPerKey, k uint32
}

func NewBloomFilter(bitsPerKey int) *BloomFilter {
	k := uint32(bitsPerKey) * 69 / 100 // 0.69 =~ ln(2)
	if k < 1 {
		k = 1
	} else if k > 30 {
		k = 30
	}
	return &BloomFilter{uint32(bitsPerKey), k}
}

func (*BloomFilter) Name() string {
	return "leveldb.BuiltinBloomFilter"
}

func (p *BloomFilter) CreateFilter(keys [][]byte, buf io.Writer) {
	bits := uint32(len(keys)) * p.bitsPerKey
	if bits < 64 {
		bits = 64
	}
	bytes := (bits + 7) / 8
	bits = bytes * 8
	array := make([]byte, bytes)
	for _, key := range keys {
		h := BloomHash(key)
		delta := (h >> 17) | (h << 15)
		for i := uint32(0); i < p.k; i++ {
			bitPos := h % bits
			array[bitPos/8] |= 1 << (bitPos % 8)
			h += delta
		}
	}
	buf.Write(array)
	buf.Write([]byte{byte(p.k)})
}

func (p *BloomFilter) KeyMayMatch(key, filter []byte) bool {
	l := uint32(len(filter))
	if l < 2 {
		return false
	}
	bits := (l - 1) * 8
	k := uint32(filter[l-1])
	if k > 30 {
		return true
	}

	h := BloomHash(key)
	delta := (h >> 17) | (h << 15)
	for i := uint32(0); i < k; i++ {
		bitPos := h % bits
		if (uint32(filter[bitPos/8]) & (1 << (bitPos % 8))) == 0 {
			return false
		}
		h += delta
	}
	return true
}
