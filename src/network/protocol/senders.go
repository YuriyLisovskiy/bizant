// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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

func (self *Protocol) SendMessage(addrTo, msgType string) bool {
	return self.sendData(addrTo, MakeRequest(msg{Type: msgType}, C_MESSAGE))
}
