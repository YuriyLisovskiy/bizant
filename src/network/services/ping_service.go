package services

import (
	"time"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/response"
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
						response.SendPing(nodeAddress, addr, knownNodes)
					}
				}
			}
		}
	}()
}
