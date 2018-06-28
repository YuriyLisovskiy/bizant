package cli

import (
	"log"
	"fmt"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockChain()
	UTXOSet := blockchain.UTXOSet{bc}
	tx := blockchain.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := blockchain.NewCoinBaseTX(from, "")
	txs := []*blockchain.Transaction{cbTx, tx}
	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	bc.CloseDB(true)
	fmt.Println("Success!")
}
