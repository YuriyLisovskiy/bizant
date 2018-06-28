package cli

import (
	"fmt"
	"strconv"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockChain("")
	bc.CloseDB(true)
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
