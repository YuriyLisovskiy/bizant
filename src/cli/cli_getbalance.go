package cli

import (
	"fmt"
	"errors"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) getBalance(address, nodeID string) error {
	if !wallet.ValidateAddress(address) {
		return errors.New(fmt.Sprintf("ERROR: Address '%s' is not valid", address))
	}
	bc := blockchain.NewBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{BlockChain: bc}
	balance := 0.0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}
	bc.CloseDB(true)
	fmt.Printf("Balance of '%s': %.6f\n", address, balance)
	return nil
}
