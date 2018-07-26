// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package tx_io

import (
	"bytes"

	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
)

type TXInput struct {
	TxId      []byte
	VOut      int
	Signature []byte
	PubKey    []byte
}

func (in TXInput) UsesKey(pubKeyHash []byte) bool {
	return bytes.Compare(wallet.HashPubKey(in.PubKey), pubKeyHash) == 0
}
