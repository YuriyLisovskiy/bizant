// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package network

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	blockchain "github.com/YuriyLisovskiy/blockchain-go/src"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/util"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/static"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/services"
	p "github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
)

func handleConnection(conn net.Conn, bc blockchain.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := util.BytesToCommand(request[:static.COMMAND_LENGTH])
	utils.PrintLog(fmt.Sprintf("Received %s command\n", command))
	switch command {
	case static.C_ADDR:
		p.HandleAddr(request)
	case static.C_BLOCK:
		p.HandleBlock(request, bc)
	case static.C_INV:
		p.HandleInv(request, bc)
	case static.C_GETBLOCKS:
		p.HandleGetBlocks(request, bc)
	case static.C_GETDATA:
		p.HandleGetData(request, bc)
	case static.C_TX:
		p.HandleTx(request, bc)
	case static.C_VERSION:
		p.HandleVersion(request, bc)
	case static.C_PING:
		p.HandlePing(request)
	case static.C_PONG:
		p.HandlePong(request)
	case static.C_ERROR:
		p.HandleError(request)
	default:
		utils.PrintLog("Unknown command!\n")
	}
	conn.Close()
}

func StartServer(nodeID, minerAddress string) {
	static.SelfNodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	if _, ok := static.KnownNodes[static.SelfNodeAddress]; ok {
		delete(static.KnownNodes, static.SelfNodeAddress)
	}
	ln, err := net.Listen(static.PROTOCOL, static.SelfNodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := blockchain.NewBlockChain(nodeID)
	pingService := &services.PingService{}
	pingService.Start(static.SelfNodeAddress, &static.KnownNodes)
	if len(minerAddress) > 0 {
		miningService := &services.MiningService{MinerAddress: minerAddress}
		miningService.Start(bc, &static.KnownNodes, &static.MemPool)
	}
	go func() {
		for nodeAddr := range static.KnownNodes {
			if nodeAddr != static.SelfNodeAddress {
				if p.SendVersion(static.SelfNodeAddress, nodeAddr, bc, &static.KnownNodes) {
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
