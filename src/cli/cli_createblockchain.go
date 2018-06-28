package cli

import (
	"fmt"
	"log"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) createBlockChain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchain.CreateBlockChain(address)
	UTXOSet := blockchain.UTXOSet{bc}
	UTXOSet.Reindex()
	bc.CloseDB(true)
	fmt.Println("Done!")
}
