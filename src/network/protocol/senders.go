// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package protocol

import (
	"io"
	"net"
	"fmt"
	"log"
	"bytes"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
)

func (self *Protocol) sendData(addr string, request []byte) bool {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		delete(*self.Config.Nodes, addr)
		fmt.Printf("\nPeers %d\n", len(*self.Config.Nodes))
		return false
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
	return true
}

func (self *Protocol) SendPing(addrFrom, addrTo string) bool {
	return self.sendData(addrTo, MakeRequest(ping{AddrFrom: addrFrom}, C_PING))
}

func (self *Protocol) SendPong(addrFrom, addrTo string) bool {
	return self.sendData(addrTo, MakeRequest(pong{AddrFrom: addrFrom}, C_PONG))
}

func (self *Protocol) SendInv(addrFrom, addrTo, kind string, items [][]byte) bool {
	return self.sendData(addrTo, MakeRequest(
		inv{
			AddrFrom: addrFrom,
			Type:     kind,
			Items:    items,
		},
		C_INV,
	))
}

func (self *Protocol) SendBlock(addrFrom, addrTo string, newBlock types.Block) bool {
	return self.sendData(addrTo, MakeRequest(
		block{
			AddrFrom: addrFrom,
			Block:    newBlock.Serialize(),
		},
		C_BLOCK,
	))
}

func (self *Protocol) SendAddr(addrTo string) bool {
	nodes := addr{}
	for knownNodeAddr := range *self.Config.Nodes {
		if knownNodeAddr != addrTo {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	return self.sendData(addrTo, MakeRequest(nodes, C_ADDR))
}

func (self *Protocol) SendGetBlocks(addrFrom, addrTo string) bool {
	return self.sendData(addrTo, MakeRequest(
		getblocks{
			AddrFrom:   addrFrom,
			BestHeight: self.Config.Chain.GetBestHeight(),
		},
		C_GETBLOCKS,
	))
}

func (self *Protocol) SendGetData(addrFrom, addrTo, kind string, id []byte) bool {
	return self.sendData(addrTo, MakeRequest(
		getdata{
			AddrFrom: addrFrom,
			Type:     kind,
			ID:       id,
		},
		C_GETDATA,
	))
}

func (self *Protocol) SendTx(addrFrom, addrTo string, tnx types.Transaction) bool {
	return self.sendData(addrTo, MakeRequest(
		tx{
			AddFrom:     addrFrom,
			Transaction: tnx.Serialize(),
		},
		C_TX,
	))
}

func (self *Protocol) SendVersion(addrFrom, addrTo string) bool {
	return self.sendData(
		addrTo,
		MakeRequest(
			version{
				Version:    NODE_VERSION,
				BestHeight: self.Config.Chain.GetBestHeight(),
				AddrFrom:   addrFrom,
			},
			C_VERSION,
		),
	)
}
