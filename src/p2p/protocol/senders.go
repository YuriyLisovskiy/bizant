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
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
)

func (p *Protocol) sendData(addr string, request []byte) bool {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		delete(*p.Config.Nodes, addr)
		fmt.Printf("\nPeers %d\n", len(*p.Config.Nodes))
		return false
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
	return true
}

func (p *Protocol) SendPing(addrFrom, addrTo string) bool {
	return p.sendData(addrTo, MakeRequest(ping{AddrFrom: addrFrom}, C_PING))
}

func (p *Protocol) SendPong(addrFrom, addrTo string) bool {
	return p.sendData(addrTo, MakeRequest(pong{AddrFrom: addrFrom}, C_PONG))
}

func (p *Protocol) SendInv(addrFrom, addrTo, kind string, items [][]byte) bool {
	return p.sendData(addrTo, MakeRequest(
		inv{
			AddrFrom: addrFrom,
			Type:     kind,
			Items:    items,
		},
		C_INV,
	))
}

func (p *Protocol) SendBlock(addrFrom, addrTo string, newBlock types.Block) bool {
	return p.sendData(addrTo, MakeRequest(
		block{
			AddrFrom: addrFrom,
			Block:    newBlock.Serialize(),
		},
		C_BLOCK,
	))
}

func (p *Protocol) SendAddr(addrTo string) bool {
	nodes := addr{}
	for knownNodeAddr := range *p.Config.Nodes {
		if knownNodeAddr != addrTo {
			nodes.AddrList = append(nodes.AddrList, knownNodeAddr)
		}
	}
	return p.sendData(addrTo, MakeRequest(nodes, C_ADDR))
}

func (p *Protocol) SendGetBlocks(addrFrom, addrTo string) bool {
	return p.sendData(addrTo, MakeRequest(
		getblocks{
			AddrFrom:   addrFrom,
			BestHeight: p.Config.Chain.GetBestHeight(),
		},
		C_GETBLOCKS,
	))
}

func (p *Protocol) SendGetData(addrFrom, addrTo, kind string, id []byte) bool {
	return p.sendData(addrTo, MakeRequest(
		getdata{
			AddrFrom: addrFrom,
			Type:     kind,
			ID:       id,
		},
		C_GETDATA,
	))
}

func (p *Protocol) SendTx(addrFrom, addrTo string, tnx types.Transaction) bool {
	return p.sendData(addrTo, MakeRequest(
		tx{
			AddFrom:     addrFrom,
			Transaction: tnx.Serialize(),
		},
		C_TX,
	))
}

func (p *Protocol) SendVersion(addrFrom, addrTo string) bool {
	return p.sendData(
		addrTo,
		MakeRequest(
			version{
				Version:    NODE_VERSION,
				BestHeight: p.Config.Chain.GetBestHeight(),
				AddrFrom:   addrFrom,
			},
			C_VERSION,
		),
	)
}

func (p *Protocol) SendMessage(addrTo, msgType string) bool {
	return p.sendData(addrTo, MakeRequest(msg{Type: msgType}, C_MESSAGE))
}
