package cli

import (
	"encoding/json"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
	"fmt"
)

func (cli *CLI) printChain(nodeID string) {
	bc := blockchain.NewBlockChain(nodeID)
	bci := bc.Iterator()
	for {
		block := bci.Next()
		data, err := json.MarshalIndent(block, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	bc.CloseDB(true)
}
