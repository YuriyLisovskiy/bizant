package response

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/request"
)

func sendData(addr string, data []byte, knownNodes *map[string]bool) bool {
	conn, err := net.Dial(utils.PROTOCOL, addr)
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
	data := utils.GobEncode(request.Ping{AddrFrom: addrFrom})
	requestData := append(utils.CommandToBytes(utils.C_PING), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendPong(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := utils.GobEncode(request.Pong{AddrFrom: addrFrom})
	pongRequest := append(utils.CommandToBytes(utils.C_PONG), data...)
	return sendData(addrTo, pongRequest, knownNodes)
}

func SendInv(addrFrom, addrTo, kind string, items [][]byte, knownNodes *map[string]bool) bool {
	inventory := request.Inv{AddrFrom: addrFrom, Type: kind, Items: items}
	data := utils.GobEncode(inventory)
	requestData := append(utils.CommandToBytes(utils.C_INV), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendBlock(addrFrom, addrTo string, block blockchain.Block, knownNodes *map[string]bool) bool {
	blockData := request.Block{AddrFrom: addrFrom, Block: block.Serialize()}
	data := utils.GobEncode(blockData)
	requestData := append(utils.CommandToBytes(utils.C_BLOCK), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendAddr(addrTo string, knownNodes *map[string]bool) bool {
	nodes := request.Addr{}
	for knownNodeAddr := range *knownNodes {
		if knownNodeAddr != addrTo {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	data := utils.GobEncode(nodes)
	requestData := append(utils.CommandToBytes(utils.C_ADDR), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendGetBlocks(addrFrom, addrTo string, knownNodes *map[string]bool) bool {
	data := utils.GobEncode(request.Getblocks{AddrFrom: addrFrom})
	requestData := append(utils.CommandToBytes(utils.C_GET_BLOCKS), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendGetData(addrFrom, addrTo, kind string, id []byte, knownNodes *map[string]bool) bool {
	data := utils.GobEncode(request.Getdata{AddrFrom: addrFrom, Type: kind, ID: id})
	requestData := append(utils.CommandToBytes(utils.C_GET_DATA), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendTx(addrFrom, addrTo string, tnx blockchain.Transaction, knownNodes *map[string]bool) bool {
	txData := request.Tx{AddFrom: addrFrom, Transaction: tnx.Serialize()}
	data := utils.GobEncode(txData)
	requestData := append(utils.CommandToBytes(utils.C_TX), data...)
	return sendData(addrTo, requestData, knownNodes)
}

func SendVersion(addrFrom, addrTo string, bc blockchain.BlockChain, knownNodes *map[string]bool) bool {
	bestHeight := bc.GetBestHeight()
	data := utils.GobEncode(request.Version{Version: utils.NODE_VERSION, BestHeight: bestHeight, AddrFrom: addrFrom})
	requestData := append(utils.CommandToBytes(utils.C_VERSION), data...)
	return sendData(addrTo, requestData, knownNodes)
}
