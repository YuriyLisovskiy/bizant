package cli

import (
	"fmt"
	"encoding/json"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) printChain(nodeID string) error {
	bc := blockchain.NewBlockChain(nodeID)
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		data, err := json.MarshalIndent(block, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	bc.CloseDB(true)
	return nil
}
