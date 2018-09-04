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

package tx_io

import (
	"bytes"

	"github.com/YuriyLisovskiy/blockchain-go/src/accounts/wallet"
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

