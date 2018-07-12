package network

import (
	"log"
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	rpcUtils "github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload rpcUtils.Addr
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	for _, newNode := range payload.AddrList {
		if newNode != selfNodeAddress {
			KnownNodes[newNode] = true
		}
	}
	utils.PrintLog(fmt.Sprintf("Peers %d\n", len(KnownNodes)))
}

func handleBlock(request []byte, bc blockchain.BlockChain) {
	blockchain.InterruptMining = true
	var buff bytes.Buffer
	var payload rpcUtils.Block
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
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
		sendGetData(payload.AddrFrom, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{BlockChain: bc}
		UTXOSet.Reindex()
	}
	blockchain.InterruptMining = false
}

func handleInv(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload rpcUtils.Inv
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	utils.PrintLog(fmt.Sprintf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type))
	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)
		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}
	if payload.Type == "tx" {
		txID := payload.Items[0]

		if memPool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload rpcUtils.Getblocks
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := bc.GetBlockHashes()
	rpcUtils.SendInv(selfNodeAddress, payload.AddrFrom, "block", blocks, &KnownNodes)
}

func handleGetData(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload rpcUtils.Getdata
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		rpcUtils.SendBlock(selfNodeAddress, payload.AddrFrom, block, &KnownNodes)
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memPool[txID]
		SendTx(payload.AddrFrom, tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload rpcUtils.Tx
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
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
	var payload rpcUtils.Version
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}
	KnownNodes[payload.AddrFrom] = true
//	if !utils.NodeIsKnown(payload.AddrFrom, KnownNodes) {
//		KnownNodes = append([]string{payload.AddrFrom}, KnownNodes...)
//	}
	for address := range KnownNodes {
		if address != selfNodeAddress {
			sendAddr(address)
		}
	}
}

func handlePing(request []byte) bool {
	var buff bytes.Buffer
	var data rpcUtils.Ping
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	payload := rpcUtils.GobEncode(rpcUtils.Pong{AddrFrom: selfNodeAddress})
	pongRequest := append(rpcUtils.CommandToBytes("pong"), payload...)
	return rpcUtils.SendData(data.AddrFrom, pongRequest, &KnownNodes)
}

func handlePong(request []byte) {
	var buff bytes.Buffer
	var data rpcUtils.Pong
	buff.Write(request[rpcUtils.COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	if data.AddrFrom != selfNodeAddress {
		KnownNodes[data.AddrFrom] = true
	}
	utils.PrintLog(fmt.Sprintf("Peers %d\n", len(KnownNodes)))
}
