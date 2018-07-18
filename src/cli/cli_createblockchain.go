// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package cli

import (
	"fmt"
	"errors"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src"
)

func (cli *CLI) createBlockChain(address, nodeId string) error {
	if !wallet.ValidateAddress(address) {
		return errors.New(fmt.Sprintf("ERROR: Address '%s' is not valid", address))
	}
	bc := src.CreateBlockChain(address, nodeId)
	UTXOSet := src.UTXOSet{BlockChain: bc}
	UTXOSet.Reindex()
	bc.CloseDB(true)
	fmt.Println("Done!")
	return nil
}
