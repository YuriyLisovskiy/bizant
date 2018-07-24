// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import (
	"hash"
	"hash/crc32"
)

var crc32tab = crc32.MakeTable(crc32.Castagnoli)

func NewCRC32C() hash.Hash32 {
	return crc32.New(crc32tab)
}

var crcMaskDelta uint32 = 0xa282ead8

// Return a masked representation of crc.
func MaskCRC32(crc uint32) uint32 {
	return ((crc >> 15) | (crc << 17)) + crcMaskDelta
}

// Return the crc whose masked representation is masked_crc.
func UnmaskCRC32(maskedCrc uint32) uint32 {
	rot := maskedCrc - crcMaskDelta
	return (rot >> 17) | (rot << 15)
}
