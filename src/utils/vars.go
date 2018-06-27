package utils

import "flag"

const BlocksBucket = "blocks"

var (
	AddBlockCmd   = flag.NewFlagSet("mine", flag.ExitOnError)
	PrintChainCmd = flag.NewFlagSet("printchain", flag.ExitOnError)
)
