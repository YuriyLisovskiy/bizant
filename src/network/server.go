package network

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/services"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/response"
	netUtils "github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

func handleConnection(conn net.Conn, bc blockchain.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := netUtils.BytesToCommand(request[:netUtils.COMMAND_LENGTH])
	utils.PrintLog(fmt.Sprintf("Received %s command\n", command))
	switch command {
	case netUtils.C_ADDR:
		handleAddr(request)
	case netUtils.C_BLOCK:
		handleBlock(request, bc)
	case netUtils.C_INV:
		handleInv(request, bc)
	case netUtils.C_GET_BLOCKS:
		handleGetBlocks(request, bc)
	case netUtils.C_GET_DATA:
		handleGetData(request, bc)
	case netUtils.C_TX:
		handleTx(request, bc)
	case netUtils.C_VERSION:
		handleVersion(request, bc)
	case netUtils.C_PING:
		handlePing(request)
	case netUtils.C_PONG:
		handlePong(request)
	case netUtils.C_ERROR:
		handleError(request)
	default:
		utils.PrintLog("Unknown command!\n")
	}
	conn.Close()
}

func StartServer(nodeID, minerAddress string) {
	SelfNodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	if _, ok := KnownNodes[SelfNodeAddress]; ok {
		delete(KnownNodes, SelfNodeAddress)
	}
	ln, err := net.Listen(netUtils.PROTOCOL, SelfNodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := blockchain.NewBlockChain(nodeID)
	pingService := &services.PingService{}
	pingService.Start(SelfNodeAddress, &KnownNodes)
	if len(minerAddress) > 0 {
		miningService := &services.MiningService{MinerAddress: minerAddress}
		miningService.Start(bc, &KnownNodes, &memPool)
	}
	go func() {
		for nodeAddr := range KnownNodes {
			if nodeAddr != SelfNodeAddress {
				if response.SendVersion(SelfNodeAddress, nodeAddr, bc, &KnownNodes) {
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
