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

package services

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

type MiningService struct {
	MinerAddress string
}

func (ms *MiningService) Start(proto *protocol.Protocol, memPool *map[string]types.Transaction) {
	go func() {
		for {
			if atomic.LoadInt32(&vars.Syncing) == 0 {
				var txs []types.Transaction
				for _, tx := range *memPool {
					if proto.Config.Chain.VerifyTransaction(tx) {

						txs = append(txs, tx)

					} else {
						// TODO: send an error to transaction's author

						utils.PrintLog(fmt.Sprintf("Invalid transaction %x\n", tx.Hash))

						data, err := json.MarshalIndent(tx, "", "  ")
						if err == nil {
							fmt.Println(string(data))
						}
					}
					delete(*memPool, hex.EncodeToString(tx.Hash))
				}
				newBlock, err := proto.Config.Chain.MineBlock(ms.MinerAddress, txs)
				if err == nil {
					utils.PrintLog("New block is mined!\n")
					go func() {
						for nodeAddr := range *proto.Config.Nodes {
							if nodeAddr != ms.MinerAddress {
								proto.SendBlock(ms.MinerAddress, nodeAddr, newBlock)
							}
						}
					}()
					UTXOSet := core.UTXOSet{BlockChain: *proto.Config.Chain}
					//	UTXOSet.Reindex()
					UTXOSet.Update(newBlock)
				}
			}
		}
	}()
}
