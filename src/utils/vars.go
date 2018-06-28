package utils

import "flag"

const (
	BlocksBucket        = "blocks"
	Reward              = 50
	DBFile              = "BlockChain.db"
	GenesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

var (
	AddBlockCmd   = flag.NewFlagSet("mine", flag.ExitOnError)
	PrintChainCmd = flag.NewFlagSet("printchain", flag.ExitOnError)
)
