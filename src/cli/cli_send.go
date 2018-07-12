package cli

import (
	"fmt"
	"errors"
	w "github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	net "github.com/YuriyLisovskiy/blockchain-go/src/network"
//	"encoding/json"
)

func (cli *CLI) send(from, to string, amount, fee float64, nodeID string) error {
	if !w.ValidateAddress(from) {
		return errors.New("ERROR: Sender address is not valid")
	}
	if !w.ValidateAddress(to) {
		return errors.New("ERROR: Recipient address is not valid")
	}
	bc := blockchain.NewBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{bc}
	wallets, err := w.NewWallets(nodeID)
	if err != nil {
		return err
	}
	wallet, err := wallets.GetWallet(from)
	if err != nil {
		return err
	}
	tx := blockchain.NewUTXOTransaction(&wallet, to, amount, fee, &UTXOSet)

//	newBlock := bc.MineBlock(from, []*blockchain.Transaction{tx})
//	UTXOSet.Update(newBlock)

//	data, err := json.MarshalIndent(tx, "", "  ")
//	if err != nil {
//		return err
//	}
//	fmt.Println(string(data))

	for nodeAddr := range net.KnownNodes {
		if nodeAddr != from {
			net.SendTx(nodeAddr, tx)
		}
	}
	bc.CloseDB(true)
	fmt.Println("Success!")
	return nil
}
