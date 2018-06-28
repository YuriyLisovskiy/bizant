package cli

import (
	"fmt"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
)

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := wallet.NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)
	fmt.Printf("Your new address: %s\n", address)
}
