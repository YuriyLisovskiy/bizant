// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package static

import "github.com/YuriyLisovskiy/blockchain-go/src/core/types"

var (
	SelfNodeAddress string

	KnownNodes = map[string]bool{
		"localhost:3000": true,
		"localhost:3001": true,
	}

	BlocksInTransit [][]byte
	MemPool         = make(map[string]types.Transaction)
)
