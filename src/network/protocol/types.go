// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package protocol

import "github.com/YuriyLisovskiy/blockchain-go/src/core"

type Configuration struct {
	Chain *core.BlockChain
	Nodes *map[string]bool
}

type Protocol struct {
	Config *Configuration
}

type Header struct {

}
