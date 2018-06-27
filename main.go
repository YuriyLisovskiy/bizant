package main

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/cli"
	"github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

func main() {
	bc := blockchain.NewBlockChain("BlockChain.db")
//	defer bc.db.Close()

	newCli := cli.NewCLI(bc)
	newCli.Run()
}
