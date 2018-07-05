package network

import "github.com/YuriyLisovskiy/blockchain-go/src/blockchain"

var nodeAddress string
var miningAddress string
var KnownNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]blockchain.Transaction)
