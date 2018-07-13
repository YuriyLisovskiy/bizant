// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package rpc

import "time"

type PingService struct {}

func (ps *PingService) Start(nodeAddress string, knownNodes *map[string]bool) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				for addr := range *knownNodes {
					if addr != nodeAddress {
						SendPing(nodeAddress, addr, knownNodes)
					}
				}
			}
		}
	}()
}
