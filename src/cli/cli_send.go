package cli

import (
	"log"
	"fmt"
	"github.com/YuriyLisovskiy/blockchain-go/src/network"
	w "github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) send(from, to string, amount, fee float64, nodeID string) {
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
	tx := blockchain.NewUTXOTransaction(&wallet, to, amount, fee, &UTXOSet)

//	newBlock := bc.MineBlock(from, []*blockchain.Transaction{tx})
//	UTXOSet.Update(newBlock)

	for nodeAddr := range network.KnownNodes {
		if nodeAddr != from {
			network.SendTx(nodeAddr, tx)
		}
	}
	bc.CloseDB(true)
	fmt.Println("Success!")
}
