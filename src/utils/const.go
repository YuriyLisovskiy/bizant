// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package utils

import "math/big"

const WORD_BYTES = (32 << (uint64(^big.Word(0)) >> 63)) / 8
