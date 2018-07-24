// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import (
	"bytes"
	"encoding/binary"
)

func Hash(data []byte, seed uint32) uint32 {
	var m uint32 = 0xc6a4a793
	var r uint32 = 24
	h := seed ^ (uint32(len(data)) * m)
	buf := bytes.NewBuffer(data)
	for buf.Len() >= 4 {
		var w uint32
		binary.Read(buf, binary.LittleEndian, &w)
		h += w
		h *= m
		h ^= h >> 16
	}
	rest := buf.Bytes()
	switch len(rest) {
	default:
		panic("not reached")
	case 3:
		h += uint32(rest[2]) << 16
		fallthrough
	case 2:
		h += uint32(rest[1]) << 8
		fallthrough
	case 1:
		h += uint32(rest[0])
		h *= m
		h ^= h >> r
	case 0:
	}
	return h
}
