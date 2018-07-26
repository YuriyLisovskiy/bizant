// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package vars

import "math"

const (
	TARGET_BITS       = 16
	MINING_REWARD     = 50.0
	MIN_CURRENCY_UNIT = 0.000001
	MIN_FEE_PER_BYTE  = 20 * MIN_CURRENCY_UNIT
	UTXO_BUCKET       = "chainstate"
	MAX_NONCE         = math.MaxInt32
)
