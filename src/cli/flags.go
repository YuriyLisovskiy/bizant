package cli

import "flag"

var (
	getBalanceCmd       = flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd = flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd     = flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd    = flag.NewFlagSet("listaddresses", flag.ExitOnError)
	sendCmd             = flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd       = flag.NewFlagSet("printchain", flag.ExitOnError)
)
