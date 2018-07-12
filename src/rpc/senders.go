package rpc

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/rpc/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func sendAddr(address string) bool {
	nodes := utils.Addr{}
	for knownNodeAddr := range KnownNodes {
		if knownNodeAddr != address {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	payload := utils.GobEncode(nodes)
	request := append(utils.CommandToBytes("addr"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func sendGetBlocks(address string) bool {
	payload := utils.GobEncode(utils.Getblocks{AddrFrom: selfNodeAddress})
	request := append(utils.CommandToBytes("getblocks"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func sendGetData(address, kind string, id []byte) bool {
	payload := utils.GobEncode(utils.Getdata{AddrFrom: selfNodeAddress, Type: kind, ID: id})
	request := append(utils.CommandToBytes("getdata"), payload...)
	return utils.SendData(address, request, &KnownNodes)
}

func SendTx(addr string, tnx blockchain.Transaction) bool {
	data := utils.Tx{AddFrom: selfNodeAddress, Transaction: tnx.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("tx"), payload...)
	return utils.SendData(addr, request, &KnownNodes)
}

func sendVersion(addr string, bc blockchain.BlockChain) bool {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(utils.Version{Version: utils.NODE_VERSION, BestHeight: bestHeight, AddrFrom: selfNodeAddress})
	request := append(utils.CommandToBytes("version"), payload...)
	return utils.SendData(addr, request, &KnownNodes)
}
