// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package cli

import (
	"fmt"
	"errors"
	"github.com/YuriyLisovskiy/blockchain-go/src/net"
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
	net.StartServer(nodeID, minerAddress)
	return nil
}
