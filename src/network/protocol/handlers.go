// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package protocol

import (
	"log"
	"fmt"
	"bytes"
	"sync/atomic"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/static"
)

func (self *Protocol) HandleAddr(request []byte) {
	var buff bytes.Buffer
	payload := addr{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	for _, newNode := range payload.AddrList {
		if newNode != static.SelfNodeAddress {
			static.KnownNodes[newNode] = true
		}
	}
	utils.PrintLog(fmt.Sprintf("Peers %d\n", len(static.KnownNodes)))
}

func (self *Protocol) HandleBlock(request []byte) {
	atomic.StoreInt32(&vars.Mining, 1)
	var buff bytes.Buffer
	payload := block{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := core.DeserializeBlock(blockData)
	utils.PrintLog("Received a new block!\n")
	self.Config.Chain.AddBlock(block)
	utils.PrintLog(fmt.Sprintf("Added block %x\n", block.Hash))
	if len(static.BlocksInTransit) > 0 {
		blockHash := static.BlocksInTransit[0]
		self.SendGetData(static.SelfNodeAddress, payload.AddrFrom, C_BLOCK, blockHash)
		static.BlocksInTransit = static.BlocksInTransit[1:]
	} else {
		UTXOSet := core.UTXOSet{BlockChain: *self.Config.Chain}
		UTXOSet.Reindex()
	}
	atomic.StoreInt32(&vars.Mining, 0)
}

func (self *Protocol) HandleInv(request []byte) {
	var buff bytes.Buffer
	payload := inv{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	utils.PrintLog(fmt.Sprintf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type))
	switch payload.Type {
	case C_BLOCK:
		static.BlocksInTransit = payload.Items
		blockHash := payload.Items[0]
		self.SendGetData(static.SelfNodeAddress, payload.AddrFrom, C_BLOCK, blockHash)
		var newInTransit [][]byte
		for _, b := range static.BlocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		static.BlocksInTransit = newInTransit
	case C_TX:
		txID := payload.Items[0]
		if static.MemPool[hex.EncodeToString(txID)].Hash == nil {
			self.SendGetData(static.SelfNodeAddress, payload.AddrFrom, C_TX, txID)
		}
	default:
	}
}

func (self *Protocol) HandleGetBlocks(request []byte) {
	var buff bytes.Buffer
	payload := getblocks{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := self.Config.Chain.GetBlockHashes(payload.BestHeight)
	self.SendInv(static.SelfNodeAddress, payload.AddrFrom, C_BLOCK, blocks)
}

func (self *Protocol) HandleGetData(request []byte) {
	var buff bytes.Buffer
	payload := getdata{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	switch payload.Type {
	case C_BLOCK:
		block, err := self.Config.Chain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		self.SendBlock(static.SelfNodeAddress, payload.AddrFrom, block)
	case C_TX:
		txID := hex.EncodeToString(payload.ID)
		tx := static.MemPool[txID]
		self.SendTx(static.SelfNodeAddress, payload.AddrFrom, tx)
		// delete(mempool, txID)
	default:
	}
}

func (self *Protocol) HandleTx(request []byte) {
	var buff bytes.Buffer
	payload := tx{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := core.DeserializeTransaction(txData)
	static.MemPool[hex.EncodeToString(tx.Hash)] = tx

	if !self.Config.Chain.VerifyTransaction(tx) {
		utils.PrintLog(fmt.Sprintf("Invalid transaction %x\n", tx.Hash))
		data, err := json.MarshalIndent(tx, "", "  ")
		if err == nil {
			fmt.Println(string(data))
		}
	}

	/*
		if selfNodeAddress == KnownNodes[0] {
			for _, node := range KnownNodes {
				if node != selfNodeAddress && node != payload.AddFrom {
					sendInv(node, "tx", [][]byte{tx.ID})
				}
			}
		} else {
			if len(memPool) >= 2 && len(miningAddress) > 0 {
			MineTransactions:
				var txs []*blockchain.Transaction
				for id := range memPool {
					tx := memPool[id]
					if bc.VerifyTransaction(&tx) {
						txs = append(txs, &tx)
					}
				}
				if len(txs) == 0 {
					fmt.Println("All transactions are invalid! Waiting for new ones...")
					return
				}
			//	newBlock := bc.MineBlock(miningAddress, txs)
				UTXOSet := blockchain.UTXOSet{bc}
				UTXOSet.Reindex()
				fmt.Println("New block is mined!")
				for _, tx := range txs {
					txID := hex.EncodeToString(tx.ID)
					delete(memPool, txID)
				}
			//	for _, node := range KnownNodes {
			//		if node != selfNodeAddress {
			//			sendInv(node, "block", [][]byte{newBlock.Hash})
			//		}
			//	}
				if len(memPool) > 0 {
					goto MineTransactions
				}
			}
		}
	*/
}

func (self *Protocol) HandleVersion(request []byte) {
	var buff bytes.Buffer
	payload := version{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	myBestHeight := self.Config.Chain.GetBestHeight()
	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		self.SendGetBlocks(static.SelfNodeAddress, payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		self.SendVersion(static.SelfNodeAddress, payload.AddrFrom)
	}
	static.KnownNodes[payload.AddrFrom] = true
	//	if !utils.NodeIsKnown(payload.AddrFrom, KnownNodes) {
	//		KnownNodes = append([]string{payload.AddrFrom}, KnownNodes...)
	//	}
	for address := range static.KnownNodes {
		if address != static.SelfNodeAddress {
			self.SendAddr(address)
		}
	}
}

func (self *Protocol) HandlePing(request []byte) bool {
	var buff bytes.Buffer
	payload := ping{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	return self.SendPong(static.SelfNodeAddress, payload.AddrFrom)
}

func (self *Protocol) HandlePong(request []byte) {
	var buff bytes.Buffer
	payload := pong{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.AddrFrom != static.SelfNodeAddress {
		static.KnownNodes[payload.AddrFrom] = true
	}
	utils.PrintLog(fmt.Sprintf("Peers %d\n", len(static.KnownNodes)))
}

func (self *Protocol) HandleError(request []byte) {

	// TODO: implement protocol error handling

}
