// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package net

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
	blockchain "github.com/YuriyLisovskiy/blockchain-go/src"
	"github.com/YuriyLisovskiy/blockchain-go/src/primitives"
)

func sendData(addr string, request []byte, knownNodes *map[string]bool) bool {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		delete(*knownNodes, addr)
		fmt.Printf("\nPeers %d\n", len(*knownNodes))
		return false
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
	return true
}

func SendPing(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(ping{AddrFrom: addrFrom}, C_PING), knownNodes)
}

func SendPong(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(pong{AddrFrom: addrFrom}, C_PONG), knownNodes)
}

func SendInv(addrFrom, addrTo, kind string, items [][]byte, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(inv{AddrFrom: addrFrom, Type: kind, Items: items}, C_INV), knownNodes)
}

func SendBlock(addrFrom, addrTo string, newBlock primitives.Block, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(block{AddrFrom: addrFrom, Block: newBlock.Serialize()}, C_BLOCK), knownNodes)
}

func SendAddr(addrTo string, knownNodes *map[string]bool) bool {
	nodes := addr{}
	for knownNodeAddr := range *knownNodes {
		if knownNodeAddr != addrTo {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	return sendData(addrTo, makeRequest(nodes, C_ADDR), knownNodes)
}

func SendGetBlocks(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(getblocks{AddrFrom: addrFrom}, C_GETBLOCKS), knownNodes)
}

func SendGetData(addrFrom, addrTo, kind string, id []byte, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(getdata{AddrFrom: addrFrom, Type: kind, ID: id}, C_GETDATA), knownNodes)
}

func SendTx(addrFrom, addrTo string, tnx primitives.Transaction, knownNodes *map[string]bool) bool {
	return sendData(addrTo, makeRequest(tx{AddFrom: addrFrom, Transaction: tnx.Serialize()}, C_TX), knownNodes)
}

func SendVersion(addrFrom, addrTo string, bc blockchain.BlockChain, knownNodes *map[string]bool) bool {
	return sendData(
		addrTo,
		makeRequest(version{Version: NODE_VERSION, BestHeight: bc.GetBestHeight(), AddrFrom: addrFrom}, C_VERSION),
		knownNodes,
	)
}
