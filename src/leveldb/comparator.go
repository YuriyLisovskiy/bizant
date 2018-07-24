// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import "bytes"

type BasicComparator interface {
	Compare(a, b []byte) int
}

type Comparator interface {
	BasicComparator
	Name() string
	FindShortestSeparator(a, b []byte) []byte
	FindShortSuccessor(b []byte) []byte
}

var DefaultComparator = BytewiseComparator{}

type BytewiseComparator struct{}

func (BytewiseComparator) Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

func (BytewiseComparator) Name() string {
	return "leveldb.BytewiseComparator"
}

func (BytewiseComparator) FindShortestSeparator(a, b []byte) []byte {
	i, n := 0, len(a)
	if n > len(b) {
		n = len(b)
	}
	for i < n && a[i] == b[i] {
		i++
	}
	if i >= n {
	} else if c := a[i]; c < 0xff && c+1 < b[i] {
		r := make([]byte, i+1)
		copy(r, a)
		r[i]++
		return r
	}
	return a
}

func (BytewiseComparator) FindShortSuccessor(b []byte) []byte {
	var res []byte
	for _, c := range b {
		if c != 0xff {
			res = append(res, c+1)
			return res
		}
		res = append(res, c)
	}
	return b
}
