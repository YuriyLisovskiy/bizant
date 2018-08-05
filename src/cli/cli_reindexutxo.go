// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

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
