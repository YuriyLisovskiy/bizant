package rpc

import "github.com/YuriyLisovskiy/blockchain-go/src/blockchain"

var (
	selfNodeAddress string
	KnownNodes      = map[string]bool{"localhost:3000": true, "localhost:3001": true}
	blocksInTransit = [][]byte{}
	memPool         = make(map[string]blockchain.Transaction)
)
