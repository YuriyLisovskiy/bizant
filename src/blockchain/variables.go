package blockchain

import "math"

var (
	maxNonce = math.MaxInt32
)
const (
	subsidy = 50
	targetBits = 16
	utxoBucket = "chainstate"
)
