// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package cli

import (
			"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"fmt"
)

func (cli *CLI) printChain(nodeID string) error {
	bc := core.NewBlockChain(nodeID)
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
//		data, err := json.MarshalIndent(block, "", "  ")
//		if err != nil {
//			return err
//		}
//		fmt.Println(string(data))
		fmt.Printf("\nBlock HASH: %x\n", block.Hash)
		fmt.Printf("Prev Block HASH: %x\n", block.PrevBlockHash)
	}
	bc.CloseDB(true)
	return nil
}
