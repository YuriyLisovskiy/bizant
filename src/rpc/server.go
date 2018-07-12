package rpc

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	rpcUtils "github.com/YuriyLisovskiy/blockchain-go/src/rpc/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/rpc/services"
)

func handleConnection(conn net.Conn, bc blockchain.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := rpcUtils.BytesToCommand(request[:rpcUtils.COMMAND_LENGTH])
	utils.PrintLog(fmt.Sprintf("Received %s command\n", command))
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
		utils.PrintLog("Unknown command!\n")
	}
	conn.Close()
}

func StartServer(nodeID, minerAddress string) {
	selfNodeAddress = fmt.Sprintf("localhost:%s", nodeID)

	delete(KnownNodes, selfNodeAddress)

	ln, err := net.Listen(rpcUtils.PROTOCOL, selfNodeAddress)
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
