// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package services

import (
	"fmt"
	"encoding/hex"
	"encoding/json"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	blockchain "github.com/YuriyLisovskiy/blockchain-go/src"
	"github.com/YuriyLisovskiy/blockchain-go/src/primitives"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
)

type MiningService struct {
	MinerAddress string
}

func (ms *MiningService) Start(bc blockchain.BlockChain, knownNodes *map[string]bool, memPool *map[string]primitives.Transaction) {
	go func() {
		for {
			var txs []primitives.Transaction
			for _, tx := range *memPool {
				if bc.VerifyTransaction(tx) {

					txs = append(txs, tx)

				} else {
					// TODO: send an error to transaction's author

					utils.PrintLog(fmt.Sprintf("Invalid transaction %x\n", tx.ID))

					data, err := json.MarshalIndent(tx, "", "  ")
					if err == nil {
						fmt.Println(string(data))
					}
				}
				delete(*memPool, hex.EncodeToString(tx.ID))
			}
			newBlock, err := bc.MineBlock(ms.MinerAddress, txs)
			if err == nil {
				utils.PrintLog("New block is mined!\n")
				go func() {
					for nodeAddr := range *knownNodes {
						if nodeAddr != ms.MinerAddress {
							protocol.SendBlock(ms.MinerAddress, nodeAddr, newBlock, knownNodes)
						}
					}
				}()
				UTXOSet := blockchain.UTXOSet{BlockChain: bc}
			//	UTXOSet.Reindex()
				UTXOSet.Update(newBlock)
			}
		}
	}()
}
