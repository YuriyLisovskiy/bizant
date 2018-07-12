package cli

import (
	"fmt"
	"errors"
	"github.com/YuriyLisovskiy/blockchain-go/src/rpc"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
)

func (cli *CLI) startNode(nodeID string, minerAddress string) error {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			return errors.New(fmt.Sprintf("wrong miner address %s", minerAddress))
		}
	}
	rpc.StartServer(nodeID, minerAddress)
	return nil
}
