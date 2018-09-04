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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/blockchain-go/src/accounts/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/static"
)

func (cli *CLI) send(from, to string, amount, fee float64, nodeID string) error {
	if !wallet.ValidateAddress(from) {
		return errors.New("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		return errors.New("ERROR: Recipient address is not valid")
	}
	bc := core.NewBlockChain(nodeID)
	utxoSet := core.UTXOSet{BlockChain: bc}
	wallets, err := wallet.NewWallets(nodeID)
	if err != nil {
		return err
	}
	senderWallet, err := wallets.GetWallet(from)
	if err != nil {
		return err
	}
	tx := core.NewUTXOTransaction(&senderWallet, to, amount, fee, &utxoSet)

//	newBlock := bc.MineBlock(from, []*blockchain.Transaction{tx})
//	UTXOSet.Update(newBlock)

	fmt.Printf("\n\nTX VERIFIED: %t\n\n", utxoSet.BlockChain.VerifyTransaction(tx))

	data, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	proto := protocol.Protocol{
		Config: &protocol.Configuration{
			Nodes: &static.KnownNodes,
			Chain: &bc,
		},
	}
	for nodeAddr := range static.KnownNodes {
		if nodeAddr != static.SelfNodeAddress {
			proto.SendTx(static.SelfNodeAddress, nodeAddr, tx)
		}
	}
	bc.CloseDB(true)
	fmt.Println("Success!")
	return nil
}
