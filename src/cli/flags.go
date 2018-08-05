// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package cli

import "flag"

var (
	getBalanceCmd       = flag.NewFlagSet("balance", flag.ExitOnError)
	createBlockChainCmd = flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd     = flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd    = flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd       = flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd      = flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd             = flag.NewFlagSet("send", flag.ExitOnError)
	startNodeCmd        = flag.NewFlagSet("startnode", flag.ExitOnError)
)
