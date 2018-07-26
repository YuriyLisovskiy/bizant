// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package cli

import (
	"fmt"
	"encoding/json"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
)

func (cli *CLI) printChain(nodeID string) error {
	bc := core.NewBlockChain(nodeID)
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		data, err := json.MarshalIndent(block, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	bc.CloseDB(true)
	return nil
}
