// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package vars

import "math"

const (
	TARGET_BITS       = 16
	MINING_REWARD     = 50.0
	MIN_CURRENCY_UNIT = 0.000001
	MIN_FEE_PER_BYTE  = 20 * MIN_CURRENCY_UNIT
	MAX_NONCE         = math.MaxInt32
)
