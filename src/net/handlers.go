// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package net

import (
	"log"
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func handleAddr(request []byte) {
	var buff bytes.Buffer
	payload := addr{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	for _, newNode := range payload.AddrList {
		if newNode != SelfNodeAddress {
			KnownNodes[newNode] = true
		}
	}
	utils.PrintLog(fmt.Sprintf("Peers %d\n", len(KnownNodes)))
}

func handleBlock(request []byte, bc blockchain.BlockChain) {
	blockchain.InterruptMining = true
	var buff bytes.Buffer
	payload := block{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)
	utils.PrintLog("Recevied a new block!\n")
	bc.AddBlock(block)
	utils.PrintLog(fmt.Sprintf("Added block %x\n", block.Hash))
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		SendGetData(SelfNodeAddress, payload.AddrFrom, C_BLOCK, blockHash, &KnownNodes)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{BlockChain: bc}
		UTXOSet.Reindex()
	}
	blockchain.InterruptMining = false
}

func handleInv(request []byte, bc blockchain.BlockChain) {
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
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		SendGetData(SelfNodeAddress, payload.AddrFrom, C_BLOCK, blockHash, &KnownNodes)
		var newInTransit [][]byte
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	case C_TX:
		txID := payload.Items[0]
		if memPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(SelfNodeAddress, payload.AddrFrom, C_TX, txID, &KnownNodes)
		}
	default:
	}
}

func handleGetBlocks(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := getblocks{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := bc.GetBlockHashes()
	SendInv(SelfNodeAddress, payload.AddrFrom, C_BLOCK, blocks, &KnownNodes)
}

func handleGetData(request []byte, bc blockchain.BlockChain) {
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
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		SendBlock(SelfNodeAddress, payload.AddrFrom, block, &KnownNodes)
	case C_TX:
		txID := hex.EncodeToString(payload.ID)
		tx := memPool[txID]
		SendTx(SelfNodeAddress, payload.AddrFrom, tx, &KnownNodes)
		// delete(mempool, txID)
	default:
	}
}

func handleTx(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := tx{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	memPool[hex.EncodeToString(tx.ID)] = tx
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

func handleVersion(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := version{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		SendGetBlocks(SelfNodeAddress, payload.AddrFrom, &KnownNodes)
	} else if myBestHeight > foreignerBestHeight {
		SendVersion(SelfNodeAddress, payload.AddrFrom, bc, &KnownNodes)
	}
	KnownNodes[payload.AddrFrom] = true
	//	if !utils.NodeIsKnown(payload.AddrFrom, KnownNodes) {
	//		KnownNodes = append([]string{payload.AddrFrom}, KnownNodes...)
	//	}
	for address := range KnownNodes {
		if address != SelfNodeAddress {
			SendAddr(address, &KnownNodes)
		}
	}
}

func handlePing(request []byte) bool {
	var buff bytes.Buffer
	data := ping{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	return SendPong(SelfNodeAddress, data.AddrFrom, &KnownNodes)
}

func handlePong(request []byte) {
	var buff bytes.Buffer
	data := pong{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	if data.AddrFrom != SelfNodeAddress {
		KnownNodes[data.AddrFrom] = true
	}
	utils.PrintLog(fmt.Sprintf("Peers %d\n", len(KnownNodes)))
}

func handleError(request []byte) {

	// TODO: implement protocol error handling

}
