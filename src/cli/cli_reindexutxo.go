package cli

import (
	"fmt"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) reindexUTXO(nodeID string) {
	bc := blockchain.NewBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{bc}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
