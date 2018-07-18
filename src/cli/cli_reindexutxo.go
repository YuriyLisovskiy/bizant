// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package cli

import (
	"fmt"
	bc "github.com/YuriyLisovskiy/blockchain-go/src"
)

func (cli *CLI) reindexUTXO(nodeID string) {
	chain := bc.NewBlockChain(nodeID)
	UTXOSet := bc.UTXOSet{BlockChain: chain}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
