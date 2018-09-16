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

package cli

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/blockchain-go/src/accounts/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/config"
	"github.com/YuriyLisovskiy/blockchain-go/src/p2p"
)

func (cli *CLI) startNode(minerAddress string) error {
	if !config.Exists() {
		return ErrConfigNotFound
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("Starting node %d\n", cfg.Port)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			return errors.New(fmt.Sprintf("wrong miner address %s", minerAddress))
		}
	}
	server := p2p.Server{}
	server.Start(cfg, minerAddress)
	return nil
}
