// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package rpc

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func handleConnection(conn net.Conn, bc blockchain.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := BytesToCommand(request[:COMMAND_LENGTH])
	utils.PrintLog(fmt.Sprintf("Received %s command\n", command))
	switch command {
	case C_ADDR:
		handleAddr(request)
	case C_BLOCK:
		handleBlock(request, bc)
	case C_INV:
		handleInv(request, bc)
	case C_GET_BLOCKS:
		handleGetBlocks(request, bc)
	case C_GET_DATA:
		handleGetData(request, bc)
	case C_TX:
		handleTx(request, bc)
	case C_VERSION:
		handleVersion(request, bc)
	case C_PING:
		handlePing(request)
	case C_PONG:
		handlePong(request)
	case C_ERROR:
		handleError(request)
	default:
		utils.PrintLog("Unknown command!\n")
	}
	conn.Close()
}

func StartRPCServer(nodeID, minerAddress string) {
	SelfNodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	if _, ok := KnownNodes[SelfNodeAddress]; ok {
		delete(KnownNodes, SelfNodeAddress)
	}
	ln, err := net.Listen(PROTOCOL, SelfNodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := blockchain.NewBlockChain(nodeID)
	pingService := &PingService{}
	pingService.Start(SelfNodeAddress, &KnownNodes)
	if len(minerAddress) > 0 {
		miningService := &MiningService{MinerAddress: minerAddress}
		miningService.Start(bc, &KnownNodes, &memPool)
	}
	go func() {
		for nodeAddr := range KnownNodes {
			if nodeAddr != SelfNodeAddress {
				if SendVersion(SelfNodeAddress, nodeAddr, bc, &KnownNodes) {
					break
				}
			}
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}
