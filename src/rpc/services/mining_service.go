package services

import (
	"encoding/hex"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	rpcUtils "github.com/YuriyLisovskiy/blockchain-go/src/rpc/utils"
	"fmt"
)

type MiningService struct {
	MinerAddress string
}

func (ms *MiningService) Start(bc blockchain.BlockChain, knownNodes *map[string]bool, memPool *map[string]blockchain.Transaction) {
	go func() {
		for {
			var txs []blockchain.Transaction
			for _, tx := range *memPool {
				if bc.VerifyTransaction(tx) {
					txs = append(txs, tx)
					delete(*memPool, hex.EncodeToString(tx.ID))
				} else {
					utils.PrintLog(fmt.Sprintf("Invalid transaction %x\n", tx.ID))
				}
			}
			newBlock, err := bc.MineBlock(ms.MinerAddress, txs)
			if err == nil {
				utils.PrintLog("New block is mined!\n")
				go func() {
					for nodeAddr := range *knownNodes {
						if nodeAddr != ms.MinerAddress {
							rpcUtils.SendBlock(ms.MinerAddress, nodeAddr, newBlock, knownNodes)
						}
					}
				}()
				UTXOSet := blockchain.UTXOSet{BlockChain: bc}
				UTXOSet.Reindex()
			}
		}
	}()
}
