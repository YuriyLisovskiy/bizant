package network

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

func sendAddr(address string) {
	nodes := addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := utils.GobEncode(nodes)
	request := append(utils.CommandToBytes("addr"), payload...)
	utils.SendData(address, request, &KnownNodes)
}

func sendBlock(addr string, b *blockchain.Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("block"), payload...)
	utils.SendData(addr, request, &KnownNodes)
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := utils.GobEncode(inventory)
	request := append(utils.CommandToBytes("inv"), payload...)
	utils.SendData(address, request, &KnownNodes)
}

func sendGetBlocks(address string) {
	payload := utils.GobEncode(getblocks{nodeAddress})
	request := append(utils.CommandToBytes("getblocks"), payload...)
	utils.SendData(address, request, &KnownNodes)
}

func sendGetData(address, kind string, id []byte) {
	payload := utils.GobEncode(getdata{nodeAddress, kind, id})
	request := append(utils.CommandToBytes("getdata"), payload...)
	utils.SendData(address, request, &KnownNodes)
}

func SendTx(addr string, tnx *blockchain.Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("tx"), payload...)
	utils.SendData(addr, request, &KnownNodes)
}

func sendVersion(addr string, bc *blockchain.BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(version{utils.NodeVersion, bestHeight, nodeAddress})
	request := append(utils.CommandToBytes("version"), payload...)
	utils.SendData(addr, request, &KnownNodes)
}
