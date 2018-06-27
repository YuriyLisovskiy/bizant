package main

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/cli"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func main() {
	bc := blockchain.NewBlockChain("BlockChain.db")
	bc.CloseDB()

	newCli := cli.NewCLI(bc)
	newCli.Run()
}
