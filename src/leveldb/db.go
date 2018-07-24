// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

type DB interface {
	Put(key, value []byte, o *WriteOptions) error
	Delete(key []byte, o *WriteOptions) error
	Write(updates *WriteBatch, o *WriteOptions) error
	Get(key []byte, o *ReadOptions) (value []byte, err error)
	NewIterator(o *ReadOptions) (iter *Iterator, err error)
	GetSnapshot() (snapshot *Snapshot, err error)
	ReleaseSnapshot(snapshot *Snapshot) error
	GetProperty(property string) (value string, err error)
	GetApproximateSizes(r *Range, n int) (size uint64, err error)
	CompactRange(begin, end []byte) error
}

type WriteBatch interface{}
type Snapshot interface{}
type Range interface{}
