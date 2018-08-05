// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

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
