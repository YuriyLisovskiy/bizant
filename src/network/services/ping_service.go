package services

import (
	"time"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/utils"
)

type PingService struct {}

func (ps *PingService) Run(nodeAddress string, knownNodes *[]string) {
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for {
			select {
			case <-ticker.C:
				for _, addr := range *knownNodes {
					utils.SendPing(nodeAddress, addr, knownNodes)
				}
			}
		}
	}()
}
