package cli

import (
	"fmt"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
)

func (cli *CLI) listAddresses(nodeID string) error {
	wallets, err := wallet.NewWallets(nodeID)
	if err != nil {
		return err
	}
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
	return nil
}
