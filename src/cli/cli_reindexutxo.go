// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package cli

import (
	"fmt"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
)

func (cli *CLI) reindexUTXO(nodeID string) {
	chain := core.NewBlockChain(nodeID)
	UTXOSet := core.UTXOSet{BlockChain: chain}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
