// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package protocol

import (
	"log"
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	blockchain "github.com/YuriyLisovskiy/blockchain-go/src"
	"github.com/YuriyLisovskiy/blockchain-go/src/primitives"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/static"
	"encoding/json"
)

func HandleAddr(request []byte) {
	var buff bytes.Buffer
	payload := addr{}
	buff.Write(request[static.COMMAND_LENGTH:])
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

func HandleBlock(request []byte, bc blockchain.BlockChain) {
	primitives.InterruptMining = true
	var buff bytes.Buffer
	payload := block{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := primitives.DeserializeBlock(blockData)
	utils.PrintLog("Recevied a new block!\n")
	bc.AddBlock(block)
	utils.PrintLog(fmt.Sprintf("Added block %x\n", block.Hash))
	if len(static.BlocksInTransit) > 0 {
		blockHash := static.BlocksInTransit[0]
		SendGetData(static.SelfNodeAddress, payload.AddrFrom, static.C_BLOCK, blockHash, &static.KnownNodes)
		static.BlocksInTransit = static.BlocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{BlockChain: bc}
		UTXOSet.Reindex()
	}
	primitives.InterruptMining = false
}

func HandleInv(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := inv{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	utils.PrintLog(fmt.Sprintf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type))
	switch payload.Type {
	case static.C_BLOCK:
		static.BlocksInTransit = payload.Items
		blockHash := payload.Items[0]
		SendGetData(static.SelfNodeAddress, payload.AddrFrom, static.C_BLOCK, blockHash, &static.KnownNodes)
		var newInTransit [][]byte
		for _, b := range static.BlocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		static.BlocksInTransit = newInTransit
	case static.C_TX:
		txID := payload.Items[0]
		if static.MemPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(static.SelfNodeAddress, payload.AddrFrom, static.C_TX, txID, &static.KnownNodes)
		}
	default:
	}
}

func HandleGetBlocks(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := getblocks{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := bc.GetBlockHashes(payload.BestHeight)
	SendInv(static.SelfNodeAddress, payload.AddrFrom, static.C_BLOCK, blocks, &static.KnownNodes)
}

func HandleGetData(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := getdata{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	switch payload.Type {
	case static.C_BLOCK:
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		SendBlock(static.SelfNodeAddress, payload.AddrFrom, block, &static.KnownNodes)
	case static.C_TX:
		txID := hex.EncodeToString(payload.ID)
		tx := static.MemPool[txID]
		SendTx(static.SelfNodeAddress, payload.AddrFrom, tx, &static.KnownNodes)
		// delete(mempool, txID)
	default:
	}
}

func HandleTx(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := tx{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := primitives.DeserializeTransaction(txData)
	static.MemPool[hex.EncodeToString(tx.ID)] = tx

	if !bc.VerifyTransaction(tx) {
		utils.PrintLog(fmt.Sprintf("Invalid transaction %x\n", tx.ID))
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

func HandleVersion(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	payload := version{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		SendGetBlocks(static.SelfNodeAddress, payload.AddrFrom, bc, &static.KnownNodes)
	} else if myBestHeight > foreignerBestHeight {
		SendVersion(static.SelfNodeAddress, payload.AddrFrom, bc, &static.KnownNodes)
	}
	static.KnownNodes[payload.AddrFrom] = true
	//	if !utils.NodeIsKnown(payload.AddrFrom, KnownNodes) {
	//		KnownNodes = append([]string{payload.AddrFrom}, KnownNodes...)
	//	}
	for address := range static.KnownNodes {
		if address != static.SelfNodeAddress {
			SendAddr(address, &static.KnownNodes)
		}
	}
}

func HandlePing(request []byte) bool {
	var buff bytes.Buffer
	payload := ping{}
	buff.Write(request[static.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	return SendPong(static.SelfNodeAddress, payload.AddrFrom, &static.KnownNodes)
}

func HandlePong(request []byte) {
	var buff bytes.Buffer
	payload := pong{}
	buff.Write(request[static.COMMAND_LENGTH:])
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

func HandleError(request []byte) {

	// TODO: implement protocol error handling

}
