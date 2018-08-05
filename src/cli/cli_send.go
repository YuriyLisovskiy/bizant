// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package cli

import (
	"fmt"
	"errors"
	"encoding/json"

	"github.com/YuriyLisovskiy/blockchain-go/src/core"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/static"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/protocol"
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
