// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import "io"

type FilterPolicy interface {
	Name() string
	CreateFilter(keys [][]byte, buf io.Writer)
	KeyMayMatch(key, filter []byte) bool
}
