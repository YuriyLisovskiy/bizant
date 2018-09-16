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

package p2p

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync/atomic"

	"github.com/YuriyLisovskiy/blockchain-go/src/config"
	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/p2p/protocol"
	"github.com/YuriyLisovskiy/blockchain-go/src/p2p/static"
	"github.com/YuriyLisovskiy/blockchain-go/src/services"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

type Server struct {
	protocol      protocol.Protocol
	pingService   services.PingService
	miningService services.MiningService
}

func handleConnection(conn net.Conn, proto *protocol.Protocol) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := protocol.BytesToCommand(request[:protocol.COMMAND_LENGTH])
	utils.PrintLog(fmt.Sprintf("Received %s command\n", command))
	switch command {
	case protocol.C_ADDR:
		proto.HandleAddr(request)
	case protocol.C_BLOCK:
		proto.HandleBlock(request)
	case protocol.C_INV:
		proto.HandleInv(request)
	case protocol.C_GETBLOCKS:
		proto.HandleGetBlocks(request)
	case protocol.C_GETDATA:
		proto.HandleGetData(request)
	case protocol.C_TX:
		proto.HandleTx(request)
	case protocol.C_VERSION:
		proto.HandleVersion(request)
	case protocol.C_PING:
		proto.HandlePing(request)
	case protocol.C_PONG:
		proto.HandlePong(request)
	case protocol.C_MESSAGE:
		proto.HandleMessage(request)
//	case protocol.C_ERROR:
//		proto.HandleError(request)
	default:
		utils.PrintLog("Unknown command!\n")
	}
	conn.Close()
}

func (s *Server) Start(cfg config.Config, minerAddress string) {
	static.SelfNodeAddress = fmt.Sprintf("%s:%d", cfg.Ip, cfg.Port)
	if _, ok := static.KnownNodes[static.SelfNodeAddress]; ok {
		delete(static.KnownNodes, static.SelfNodeAddress)
	}
	ln, err := net.Listen(protocol.PROTOCOL, static.SelfNodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := core.NewBlockChain(cfg)

	s.protocol = protocol.Protocol{
		Config: &protocol.Configuration{
			Chain: &bc,
			Nodes: &static.KnownNodes,
		},
	}
	pingService := &services.PingService{}
	pingService.Start(static.SelfNodeAddress, &s.protocol)
	go s.SyncDB()
	go func() {
		if len(minerAddress) > 0 {
			miningService := &services.MiningService{MinerAddress: minerAddress}
			miningService.Start(&s.protocol, &static.MemPool)
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, &s.protocol)
	}
}

func (s *Server) SyncDB() {
	atomic.StoreInt32(&vars.Syncing, 1)
	for nodeAddr := range static.KnownNodes {
		if nodeAddr != static.SelfNodeAddress {
			if s.protocol.SendVersion(static.SelfNodeAddress, nodeAddr) {
				break
			}
		}
	}
	if len(static.KnownNodes) < 1 {
		atomic.StoreInt32(&vars.Syncing, 0)
	}
}
