package network

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

func sendAddr(address string) bool {
	nodes := addr{}
	for knownNodeAddr := range KnownNodes {
		if knownNodeAddr != address {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	payload := utils.GobEncode(nodes)
	request := append(utils.CommandToBytes("addr"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func sendBlock(addr string, b *blockchain.Block) bool {
	data := block{selfNodeAddress, b.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("block"), payload...)
	return utils.SendData(addr, request, &KnownNodes)
}

func sendInv(address, kind string, items [][]byte) bool {
	inventory := inv{selfNodeAddress, kind, items}
	payload := utils.GobEncode(inventory)
	request := append(utils.CommandToBytes("inv"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func sendGetBlocks(address string) bool {
	payload := utils.GobEncode(getblocks{selfNodeAddress})
	request := append(utils.CommandToBytes("getblocks"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func sendGetData(address, kind string, id []byte) bool {
	payload := utils.GobEncode(getdata{selfNodeAddress, kind, id})
	request := append(utils.CommandToBytes("getdata"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func SendTx(addr string, tnx *blockchain.Transaction) bool {
	data := tx{selfNodeAddress, tnx.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("tx"), payload...)
	return utils.SendData(addr, request, &KnownNodes)
}

func sendVersion(addr string, bc *blockchain.BlockChain) bool {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(version{utils.NODE_VERSION, bestHeight, selfNodeAddress})
	request := append(utils.CommandToBytes("version"), payload...)
	return utils.SendData(addr, request, &KnownNodes)
}
