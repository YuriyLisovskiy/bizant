// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package cli

import "flag"

var (
	getBalanceCmd       = flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd = flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd     = flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd    = flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd       = flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd      = flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd             = flag.NewFlagSet("send", flag.ExitOnError)
	startNodeCmd        = flag.NewFlagSet("startnode", flag.ExitOnError)
)
