// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package services

import (
	"fmt"
	"sync/atomic"
	"encoding/hex"
	"encoding/json"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
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
