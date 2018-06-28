package cli

import (
	"fmt"
	"log"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchain.NewBlockChain()
	UTXOSet := blockchain.UTXOSet{bc}
	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}
	bc.CloseDB(true)
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
