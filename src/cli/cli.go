package cli

import (
	"os"
	"fmt"
	"strconv"
	cliVars "github.com/YuriyLisovskiy/blockchain-go/src/utils"
	chain "github.com/YuriyLisovskiy/blockchain-go/src/blockchain"
)

type CLI struct {
	bc *chain.BlockChain
}

func NewCLI(bc *chain.BlockChain) CLI {
	return CLI{bc}
}

func (cli *CLI) Run() {

//	TODO: validate args

	addBlockData := cliVars.AddBlockCmd.String("data", "", "Block data")
	if len(os.Args) == 1 {
		cliVars.AddBlockCmd.Usage()
		return
	}
	switch os.Args[1] {
	case "mine":
		err := cliVars.AddBlockCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "printchain":
		err := cliVars.PrintChainCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	default:
		cli.printUsage()
		return
	}
	if cliVars.AddBlockCmd.Parsed() {
		if *addBlockData == "" {
			cliVars.AddBlockCmd.Usage()
			return
		}
		cli.addBlock(*addBlockData)
	}
	if cliVars.PrintChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printUsage() {
	cliVars.AddBlockCmd.Usage()
	cliVars.PrintChainCmd.Usage()
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := chain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
