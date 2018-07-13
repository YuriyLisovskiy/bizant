package cli

import (
	"fmt"
	"errors"
	w "github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"github.com/YuriyLisovskiy/blockchain-go/src/network/response"
	net "github.com/YuriyLisovskiy/blockchain-go/src/network"
	"encoding/json"
)

func (cli *CLI) send(from, to string, amount, fee float64, nodeID string) error {
	if !w.ValidateAddress(from) {
		return errors.New("ERROR: Sender address is not valid")
	}
	if !w.ValidateAddress(to) {
		return errors.New("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockChain(nodeID)
	utxoSet := blockchain.UTXOSet{BlockChain: bc}
	wallets, err := w.NewWallets(nodeID)
	if err != nil {
		return err
	}
	wallet, err := wallets.GetWallet(from)
	if err != nil {
		return err
	}
	tx := blockchain.NewUTXOTransaction(&wallet, to, amount, fee, &utxoSet)

//	newBlock := bc.MineBlock(from, []*blockchain.Transaction{tx})
//	UTXOSet.Update(newBlock)

	fmt.Printf("\n\nTX VERIFIED: %t\n\n", utxoSet.BlockChain.VerifyTransaction(tx))

	data, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	for nodeAddr := range net.KnownNodes {
		if nodeAddr != net.SelfNodeAddress {
			response.SendTx(net.SelfNodeAddress, nodeAddr, tx, &net.KnownNodes)
		}
	}
	bc.CloseDB(true)
	fmt.Println("Success!")
	return nil
}
