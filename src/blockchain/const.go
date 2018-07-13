// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package blockchain

const (
	MINING_REWARD     = 50.0
	TARGET_BITS       = 16
	UTXO_BUCKET       = "chainstate"
	MIN_CURRENCY_UNIT = 0.000001
	MIN_FEE_PER_BYTE  = 20 * MIN_CURRENCY_UNIT
)
