// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package network

import "github.com/YuriyLisovskiy/blockchain-go/src/primitives"

var (
	SelfNodeAddress string

	KnownNodes = map[string]bool{
		"localhost:3000": true,
		"localhost:3001": true,
	}

	blocksInTransit [][]byte
	memPool         = make(map[string]primitives.Transaction)
)
