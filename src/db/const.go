// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

const Version = 1

const (
	MaxKeySize  = 511
	MaxDataSize = 0xffffffff
)

const (
	DefaultMapSize     = 1048576
	DefaultReaderCount = 126
)
