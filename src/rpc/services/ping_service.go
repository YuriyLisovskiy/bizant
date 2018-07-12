package services

import (
	"time"
	"github.com/YuriyLisovskiy/blockchain-go/src/rpc/utils"
)

type PingService struct {}

func (ps *PingService) Start(nodeAddress string, knownNodes *map[string]bool) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				for addr := range *knownNodes {
					if addr != nodeAddress {
						utils.SendPing(nodeAddress, addr, knownNodes)
					}
				}
			}
		}
	}()
}
