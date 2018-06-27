package main

import (
	"fmt"
	"github.com/YuriyLisovskiy/BlockChainGo/src"
	"strconv"
)

func main() {
	bc := src.NewBlockChain("BlockChain.db")

//	bc.AddBlock("Send 1 BTC to Ivan")
//	bc.AddBlock("Send 2 more BTC to Ivan")

	iterator := bc.Iterator()

	for {
		block := iterator.Next()
		if block == nil {
			break
		}
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := src.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
