// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

type Cache interface {
	Insert(key, value []byte, charge int) CacheObject
	Lookup(key []byte) CacheObject
	Erase(key []byte)
	NewId() uint64
}

type CacheObject interface {
	Release()
	Value() []byte
}
