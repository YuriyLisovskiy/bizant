// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package services

import (
	"time"

	"github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
)

type PingService struct {}

func (ps *PingService) Start(nodeAddress string, proto *protocol.Protocol) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				for addr := range *proto.Config.Nodes {
					if addr != nodeAddress {
						proto.SendPing(nodeAddress, addr)
					}
				}
			}
		}
	}()
}
