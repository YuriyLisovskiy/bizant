// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package tx_io

import (
	"bytes"

	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
)

type TXInput struct {
	PreviousTx []byte
	VOut       int
	Signature  []byte
	PubKey    []byte
}

func (in TXInput) UsesKey(pubKeyHash []byte) bool {
	return bytes.Compare(wallet.HashPubKey(in.PubKey), pubKeyHash) == 0
}

