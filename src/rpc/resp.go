// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package rpc

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func sendData(addr string, data []byte, knownNodes *map[string]bool) bool {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		delete(*knownNodes, addr)
		fmt.Printf("\nPeers %d\n", len(*knownNodes))
		return false
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
	return true
}

func SendPing(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := GobEncode(ping{AddrFrom: addrFrom})
	requestData := append(CommandToBytes(C_PING), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendPong(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := GobEncode(pong{AddrFrom: addrFrom})
	pongRequest := append(CommandToBytes(C_PONG), data...)
	return sendData(addrTo, pongRequest, knownNodes)
}

func SendInv(addrFrom, addrTo, kind string, items [][]byte, knownNodes *map[string]bool) bool {
	inventory := inv{AddrFrom: addrFrom, Type: kind, Items: items}
	data := GobEncode(inventory)
	requestData := append(CommandToBytes(C_INV), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendBlock(addrFrom, addrTo string, newBlock blockchain.Block, knownNodes *map[string]bool) bool {
	blockData := block{AddrFrom: addrFrom, Block: newBlock.Serialize()}
	data := GobEncode(blockData)
	requestData := append(CommandToBytes(C_BLOCK), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendAddr(addrTo string, knownNodes *map[string]bool) bool {
	nodes := addr{}
	for knownNodeAddr := range *knownNodes {
		if knownNodeAddr != addrTo {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	data := GobEncode(nodes)
	requestData := append(CommandToBytes(C_ADDR), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendGetBlocks(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := GobEncode(getblocks{AddrFrom: addrFrom})
	requestData := append(CommandToBytes(C_GET_BLOCKS), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendGetData(addrFrom, addrTo, kind string, id []byte, knownNodes *map[string]bool) bool {
	data := GobEncode(getdata{AddrFrom: addrFrom, Type: kind, ID: id})
	requestData := append(CommandToBytes(C_GET_DATA), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendTx(addrFrom, addrTo string, tnx blockchain.Transaction, knownNodes *map[string]bool) bool {
	txData := tx{AddFrom: addrFrom, Transaction: tnx.Serialize()}
	data := GobEncode(txData)
	requestData := append(CommandToBytes(C_TX), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendVersion(addrFrom, addrTo string, bc blockchain.BlockChain, knownNodes *map[string]bool) bool {
	bestHeight := bc.GetBestHeight()
	data := GobEncode(version{Version: NODE_VERSION, BestHeight: bestHeight, AddrFrom: addrFrom})
	requestData := append(CommandToBytes(C_VERSION), data...)
	return sendData(addrTo, requestData, knownNodes)
}
