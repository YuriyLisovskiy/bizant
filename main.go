package main

import (
	"github.com/YuriyLisovskiy/BlockChainGo/src/cli"
	"github.com/YuriyLisovskiy/BlockChainGo/src/blockchain"
)

func main() {
	bc := blockchain.NewBlockChain("BlockChain.db")
//	defer bc.db.Close()

	newCli := cli.NewCLI(bc)
	newCli.Run()
}
