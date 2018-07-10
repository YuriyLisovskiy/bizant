package network

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	gUtils "github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/services"
)

func handleConnection(conn net.Conn, bc *blockchain.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := utils.BytesToCommand(request[:utils.COMMAND_LENGTH])
	gUtils.PrintLog(fmt.Sprintf("Received %s command\n", command))
	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	case "ping":
		handlePing(request)
	case "pong":
		handlePong(request)
	default:
		gUtils.PrintLog("Unknown command!\n")
	}
	conn.Close()
}

func StartServer(nodeID, minerAddress string) {
	selfNodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	ln, err := net.Listen(utils.PROTOCOL, selfNodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := blockchain.NewBlockChain(nodeID)
	pingService := &services.PingService{}
	pingService.Start(selfNodeAddress, &KnownNodes)
	if len(minerAddress) > 0 {
		miningService := &services.MiningService{MinerAddress: minerAddress}
		miningService.Start(bc, &KnownNodes, &memPool)
	}
	go func() {
		for nodeAddr := range KnownNodes {
			if nodeAddr != selfNodeAddress {
				if sendVersion(nodeAddr, bc) {
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
