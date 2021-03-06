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
	"fmt"

	"github.com/YuriyLisovskiy/blockchain-go/src/accounts/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/config"
)

func (cli *CLI) listAddresses(cfg config.Config) error {
	wallets, err := wallet.NewWallets(cfg)
	if err != nil {
		return err
	}
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
	return nil
}
