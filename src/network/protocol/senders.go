// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package protocol

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
	blockchain "github.com/YuriyLisovskiy/blockchain-go/src"
	"github.com/YuriyLisovskiy/blockchain-go/src/primitives"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/util"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/static"
)

func sendData(addr string, request []byte, knownNodes *map[string]bool) bool {
	conn, err := net.Dial(static.PROTOCOL, addr)
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
	return sendData(addrTo, util.MakeRequest(ping{AddrFrom: addrFrom}, static.C_PING), knownNodes)
}

func SendPong(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	return sendData(addrTo, util.MakeRequest(pong{AddrFrom: addrFrom}, static.C_PONG), knownNodes)
}

func SendInv(addrFrom, addrTo, kind string, items [][]byte, knownNodes *map[string]bool) bool {
	return sendData(addrTo, util.MakeRequest(inv{AddrFrom: addrFrom, Type: kind, Items: items}, static.C_INV), knownNodes)
}

func SendBlock(addrFrom, addrTo string, newBlock primitives.Block, knownNodes *map[string]bool) bool {
	return sendData(addrTo, util.MakeRequest(block{AddrFrom: addrFrom, Block: newBlock.Serialize()}, static.C_BLOCK), knownNodes)
}

func SendAddr(addrTo string, knownNodes *map[string]bool) bool {
	nodes := addr{}
	for knownNodeAddr := range *knownNodes {
		if knownNodeAddr != addrTo {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	return sendData(addrTo, util.MakeRequest(nodes, static.C_ADDR), knownNodes)
}

func SendGetBlocks(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	return sendData(addrTo, util.MakeRequest(getblocks{AddrFrom: addrFrom}, static.C_GETBLOCKS), knownNodes)
}

func SendGetData(addrFrom, addrTo, kind string, id []byte, knownNodes *map[string]bool) bool {
	return sendData(addrTo, util.MakeRequest(getdata{AddrFrom: addrFrom, Type: kind, ID: id}, static.C_GETDATA), knownNodes)
}

func SendTx(addrFrom, addrTo string, tnx primitives.Transaction, knownNodes *map[string]bool) bool {
	return sendData(addrTo, util.MakeRequest(tx{AddFrom: addrFrom, Transaction: tnx.Serialize()}, static.C_TX), knownNodes)
}

func SendVersion(addrFrom, addrTo string, bc blockchain.BlockChain, knownNodes *map[string]bool) bool {
	return sendData(
		addrTo,
		util.MakeRequest(version{
			Version: static.NODE_VERSION,
			BestHeight: bc.GetBestHeight(),
			AddrFrom: addrFrom,
		}, static.C_VERSION),
		knownNodes,
	)
}
