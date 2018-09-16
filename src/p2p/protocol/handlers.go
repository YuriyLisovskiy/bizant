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

package protocol

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/p2p/static"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

func (*Protocol) HandleAddr(request []byte) {
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

func (p *Protocol) HandleBlock(request []byte) {
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
	p.Config.Chain.AddBlock(block)
	utils.PrintLog(fmt.Sprintf("Added block %x\n", block.Hash))
	if len(static.BlocksInTransit) > 0 {
		blockHash := static.BlocksInTransit[0]
		p.SendGetData(static.SelfNodeAddress, payload.AddrFrom, C_BLOCK, blockHash)
		static.BlocksInTransit = static.BlocksInTransit[1:]
	} else {
		UTXOSet := core.UTXOSet{BlockChain: *p.Config.Chain}
		UTXOSet.Reindex()
		atomic.StoreInt32(&vars.Syncing, 0)
	}
}

func (p *Protocol) HandleInv(request []byte) {
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
		var newInTransit [][]byte
		for _, b := range static.BlocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		static.BlocksInTransit = newInTransit
		p.SendGetData(static.SelfNodeAddress, payload.AddrFrom, C_BLOCK, blockHash)
	case C_TX:
		txID := payload.Items[0]
		if static.MemPool[hex.EncodeToString(txID)].Hash == nil {
			p.SendGetData(static.SelfNodeAddress, payload.AddrFrom, C_TX, txID)
		}
	default:
	}
}

func (p *Protocol) HandleGetBlocks(request []byte) {
	var buff bytes.Buffer
	payload := getblocks{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := p.Config.Chain.GetBlockHashes(payload.BestHeight)
	p.SendInv(static.SelfNodeAddress, payload.AddrFrom, C_BLOCK, blocks)
}

func (p *Protocol) HandleGetData(request []byte) {
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
		block, err := p.Config.Chain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		p.SendBlock(static.SelfNodeAddress, payload.AddrFrom, block)
	case C_TX:
		txID := hex.EncodeToString(payload.ID)
		tx := static.MemPool[txID]
		p.SendTx(static.SelfNodeAddress, payload.AddrFrom, tx)
		// delete(mempool, txID)
	default:
	}
}

func (p *Protocol) HandleTx(request []byte) {
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

	if !p.Config.Chain.VerifyTransaction(tx) {
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

func (p *Protocol) HandleVersion(request []byte) {
	var buff bytes.Buffer
	payload := version{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	myBestHeight := p.Config.Chain.GetBestHeight()
	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		atomic.StoreInt32(&vars.Syncing, 1)
		p.SendGetBlocks(static.SelfNodeAddress, payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		p.SendVersion(static.SelfNodeAddress, payload.AddrFrom)
	} else {
		p.SendMessage(payload.AddrFrom, C_SYNCED)
		atomic.StoreInt32(&vars.Syncing, 0)
	}
	static.KnownNodes[payload.AddrFrom] = true
	//	if !utils.NodeIsKnown(payload.AddrFrom, KnownNodes) {
	//		KnownNodes = append([]string{payload.AddrFrom}, KnownNodes...)
	//	}
	for address := range static.KnownNodes {
		if address != static.SelfNodeAddress {
			p.SendAddr(address)
		}
	}
}

func (p *Protocol) HandlePing(request []byte) bool {
	var buff bytes.Buffer
	payload := ping{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	return p.SendPong(static.SelfNodeAddress, payload.AddrFrom)
}

func (*Protocol) HandlePong(request []byte) {
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

func (*Protocol) HandleMessage(request []byte) {
	var buff bytes.Buffer
	payload := msg{}
	buff.Write(request[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	switch payload.Type {
	case C_SYNCED:
		atomic.StoreInt32(&vars.Syncing, 0)
	default:
		utils.PrintLog("Unknown msg type!\n")
	}
}

func (*Protocol) HandleError(request []byte) {

	// TODO: implement protocol error handling

}
