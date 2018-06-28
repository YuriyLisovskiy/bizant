package cli

import (
	"log"
	"fmt"
	w "github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !w.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !w.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{bc}
	wallets, err := w.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := blockchain.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := blockchain.NewCoinBaseTX(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
	}
	bc.CloseDB(true)
	fmt.Println("Success!")
}
