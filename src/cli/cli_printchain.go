// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"fmt"

	"github.com/YuriyLisovskiy/blockchain-go/src/config"
	"github.com/YuriyLisovskiy/blockchain-go/src/core"
)

func (cli *CLI) printChain(cfg config.Config) error {
	bc := core.NewBlockChain(cfg)
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
