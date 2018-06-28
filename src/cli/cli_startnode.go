package cli

import (
	"fmt"
	"log"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/network"
)

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(nodeID, minerAddress)
}
