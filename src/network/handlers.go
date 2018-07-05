package network

import (
	"log"
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

func requestBlocks() {
	for _, node := range KnownNodes {
		sendGetBlocks(node)
	}
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(KnownNodes))
	requestBlocks()
}

func handleBlock(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload block
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)
	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)
	fmt.Printf("Added block %x\n", block.Hash)
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

func handleInv(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload inv
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)
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

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload getblocks
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload getdata
	buff.Write(request[utils.CommandLength:])
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
		sendBlock(payload.AddrFrom, &block)
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]
		SendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload tx
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx
	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*blockchain.Transaction
			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}
			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}
			cbTx := blockchain.NewCoinBaseTX(miningAddress, "")
			txs = append(txs, cbTx)
			newBlock := bc.MineBlock(txs)
			UTXOSet := blockchain.UTXOSet{bc}
			UTXOSet.Reindex()
			fmt.Println("New block is mined!")
			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}
			for _, node := range KnownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}
			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload version
	buff.Write(request[utils.CommandLength:])
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

	//	sendAddr(payload.AddrFrom)

	if !nodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

func handlePing(request []byte) {
	var buff bytes.Buffer
	var requestData utils.Ping
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&requestData)
	if err != nil {
		log.Panic(err)
	}
	payload := utils.GobEncode(utils.Pong{AddrFrom: nodeAddress})
	pongRequest := append(utils.CommandToBytes("pong"), payload...)
	utils.SendData(requestData.AddrFrom, pongRequest, &KnownNodes)
}

func handlePong(request []byte) {
	var buff bytes.Buffer
	var data utils.Pong
	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	if !nodeIsKnown(data.AddrFrom) {
		KnownNodes = append(KnownNodes, data.AddrFrom)
	}
	fmt.Printf("Peers %d\n", len(KnownNodes))
}
