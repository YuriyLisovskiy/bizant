// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/YuriyLisovskiy/blockchain-go/src/config"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Print("  config\n    -ip string\n\tNode ip address\n    -port\n\tNode id\n    -path.chain\n\tPath to block chain database\n    -path.wallets\n\tPath to wallets location\n    -default\n\tSet default config\n\n")
	fmt.Print("  createblockchain\n    -address string\n\tThe address to send genesis block reward to\n\n")
	fmt.Print("  createwallet\n\tGenerates a new key-pair and saves it into the wallet file\n\n")
	fmt.Print("  getbalance\n    -address string\n\tThe address to get balance for\n\n")
	fmt.Print("  listaddresses\n\tLists all addresses from the wallet file\n\n")
	fmt.Print("  printchain\n\tPrint all the blocks of the blockchain\n\n")
	fmt.Print("  reindexutxo\n\tRebuilds the UTXO set\n\n")
	fmt.Print("  send\n    -from string\n\tSource wallet address\n    -to string\n\tDestination wallet address\n    -amount int\n\tAmount to send\n    -mine\n\tMine on the same node\n\n")
	fmt.Print("  startnode\n    -miner string\n\tStart a node with ID specified in NODE_ID env. var. -miner enables mining\n\n")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")

	configIp := configCmd.String("ip", "", "Node ip address")
	configPort := configCmd.Int("port", -1, "Node id")
	configChainPath := configCmd.String("path.chain", "", "Path to block chain database")
	configWalletsPath := configCmd.String("path.wallets", "", "Path to wallets location")
	configDefault := configCmd.Bool("default", false, "Set default config")

	createBlockChainAddress := createBlockChainCmd.String("address", "", "The address to send genesis block reward to")

	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Float64("amount", 0, "Amount to send")
	sendFee := sendCmd.Float64("fee", vars.MIN_FEE_PER_BYTE, "Mine immediately on the same node")

	startNodeMiner := startNodeCmd.String("mine", "", "Enable mining mode")

	switch os.Args[1] {
	case "balance":
		checkError(getBalanceCmd.Parse(os.Args[2:]))
	case "config":
		checkError(configCmd.Parse(os.Args[2:]))
	case "createblockchain":
		checkError(createBlockChainCmd.Parse(os.Args[2:]))
	case "createwallet":
		checkError(createWalletCmd.Parse(os.Args[2:]))
	case "listaddresses":
		checkError(listAddressesCmd.Parse(os.Args[2:]))
	case "printchain":
		checkError(printChainCmd.Parse(os.Args[2:]))
	case "reindexutxo":
		checkError(reindexUTXOCmd.Parse(os.Args[2:]))
	case "send":
		checkError(sendCmd.Parse(os.Args[2:]))
	case "startnode":
		checkError(startNodeCmd.Parse(os.Args[2:]))
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if configCmd.Parsed() {
		if *configDefault {
			cli.setDefaultConfig()
		} else {
			cli.setConfig(*configIp, *configPort, *configChainPath, *configWalletsPath)
		}
	}
	if !config.Exists() {
		log.Println(ErrConfigNotFound)
		return
	}
	cfg, err := config.LoadConfig()
	checkError(err)
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		checkError(cli.getBalance(*getBalanceAddress, cfg))
	}
	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
		checkError(cli.createBlockChain(*createBlockChainAddress, cfg))
	}
	if createWalletCmd.Parsed() {
		cli.createWallet(cfg)
	}
	if listAddressesCmd.Parsed() {
		checkError(cli.listAddresses(cfg))
	}
	if printChainCmd.Parsed() {
		checkError(cli.printChain(cfg))
	}
	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(cfg)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		checkError(cli.send(*sendFrom, *sendTo, *sendAmount, *sendFee, cfg))
	}
	if startNodeCmd.Parsed() {
		checkError(cli.startNode(*startNodeMiner))
	}
}
