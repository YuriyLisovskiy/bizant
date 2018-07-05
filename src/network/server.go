package network

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/services"
)

func handleConnection(conn net.Conn, bc *blockchain.BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := utils.BytesToCommand(request[:utils.COMMAND_LENGTH])
	fmt.Printf("Received %s command\n", command)
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
		fmt.Println("Unknown command!")
	}
	conn.Close()
}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(utils.PROTOCOL, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := blockchain.NewBlockChain(nodeID)
	if nodeAddress != KnownNodes[0] {
		sendVersion(KnownNodes[0], bc)
	}
	pingService := &services.PingService{}
	pingService.Start(nodeAddress, &KnownNodes)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}
