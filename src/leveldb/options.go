// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

// Database flag
type OptionsFlag uint

const (
	OFCreateIfMissing OptionsFlag = 1 << iota
	OFErrorIfExist
	OFParanoidCheck
)

// Database compression type
type Compression uint

const (
	DefaultCompression Compression = iota
	NoCompression
	SnappyCompression
	nCompression
)

// Database options
type Options struct {
	Comparator           Comparator
	Flag                 OptionsFlag
	Env                  Env
	WriteBuffer          int
	MaxOpenFiles         int
	BlockCache           Cache
	BlockSize            int
	BlockRestartInterval int
	CompressionType      Compression
	FilterPolicy         FilterPolicy
}

func (o *Options) GetComparator() Comparator {
	if o == nil || o.Comparator == nil {
		return DefaultComparator
	}
	return o.Comparator
}

func (o *Options) HasFlag(flag OptionsFlag) bool {
	if o == nil {
		return false
	}
	return (o.Flag | flag) != 0
}

func (o *Options) GetEnv() Env {
	if o == nil || o.Env == nil {
		return DefaultEnv
	}
	return o.Env
}

func (o *Options) GetWriteBuffer() int {
	if o == nil || o.WriteBuffer <= 0 {
		return 4 << 20
	}
	return o.WriteBuffer
}

func (o *Options) GetMaxOpenFiles() int {
	if o == nil || o.MaxOpenFiles <= 0 {
		return 1000
	}
	return o.MaxOpenFiles
}

func (o *Options) GetBlockCache() Cache {
	if o == nil {
		return nil
	}
	return o.BlockCache
}

func (o *Options) GetBlockSize() int {
	if o == nil || o.BlockSize <= 0 {
		return 4096
	}
	return o.BlockSize
}

func (o *Options) GetBlockRestartInterval() int {
	if o == nil || o.BlockRestartInterval <= 0 {
		return 16
	}
	return o.BlockRestartInterval
}

func (o *Options) GetCompressionType() Compression {
	if o == nil || o.CompressionType <= DefaultCompression || o.CompressionType >= nCompression {
		return SnappyCompression
	}
	return o.CompressionType
}

func (o *Options) GetFilterPolicy() FilterPolicy {
	if o == nil {
		return nil
	}
	return o.FilterPolicy
}

// Read flag
type ReadOptionsFlag uint

const (
	RFVerifyChecksums ReadOptionsFlag = 1 << iota
	RFDontFillCache
)

// Read options
type ReadOptions struct {
	// Specify the read flag
	Flag ReadOptionsFlag

	// If "snapshot" is non-NULL, read as of the supplied snapshot
	// (which must belong to the DB that is being read and which must
	// not have been released).  If "snapshot" is NULL, use an impliicit
	// snapshot of the state at the beginning of this read operation.
	// Default: NULL
	Snapshot Snapshot
}

func (o *ReadOptions) HasFlag(flag ReadOptionsFlag) bool {
	if o == nil {
		return false
	}
	return (o.Flag | flag) != 0
}

func (o *ReadOptions) GetSnapshot() Snapshot {
	if o == nil {
		return nil
	}
	return o.Snapshot
}

// Write flag
type WriteOptionsFlag uint

const (
	WFSync WriteOptionsFlag = 1 << iota
)

// Write options
type WriteOptions struct {
	Flag WriteOptionsFlag
}

func (o *WriteOptions) HasFlag(flag WriteOptionsFlag) bool {
	if o == nil {
		return false
	}
	return (o.Flag | flag) != 0
}
