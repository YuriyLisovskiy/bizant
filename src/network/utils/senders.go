package utils

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

type Ping struct {
	AddrFrom string
}

type Pong struct {
	AddrFrom string
}

func SendData(addr string, data []byte, knownNodes *map[string]bool) bool {
	conn, err := net.Dial(PROTOCOL, addr)
	success := true
	if err != nil {
		success = false
		delete(*knownNodes, addr)
		fmt.Printf("\nPeers %d\n", len(*knownNodes))
		return success
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
	return success
}

func SendPing(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := GobEncode(Ping{AddrFrom: addrFrom})
	request := append(CommandToBytes("ping"), data...)
	return SendData(addrTo, request, knownNodes)
}

func SendInv(selfAddress, address, kind string, items [][]byte, knownNodes *map[string]bool) bool {
	inventory := Inv{selfAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CommandToBytes("inv"), payload...)
	return SendData(address, request, knownNodes)
}

func SendBlock(addrFrom, addrTo string, block *blockchain.Block, knownNodes *map[string]bool) bool {

	println("Height:", block.Height)

	data := Block{AddrFrom: addrFrom, Block: block.Serialize()}
	payload := GobEncode(data)
	request := append(CommandToBytes("block"), payload...)
	return SendData(addrTo, request, knownNodes)
}
