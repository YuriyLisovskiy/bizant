// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package utils

import "math/big"

const (
	WORD_BYTES    = (32 << (uint64(^big.Word(0)) >> 63)) / 8
	BLOCKS_BUCKET = "blocks"
)
