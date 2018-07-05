package network

import "github.com/YuriyLisovskiy/blockchain-go/src/blockchain"

var (
	nodeAddress     string
	miningAddress   string
	KnownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	memPool         = make(map[string]blockchain.Transaction)
)
