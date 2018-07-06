package services

import (
	"time"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

type PingService struct {}

func (ps *PingService) Start(nodeAddress string, knownNodes *[]string) {
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for {
			select {
			case <-ticker.C:
				for _, addr := range *knownNodes {
					if addr != nodeAddress {
						utils.SendPing(nodeAddress, addr, knownNodes)
					}
				}
			}
		}
	}()
}
