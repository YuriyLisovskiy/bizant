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

package core

import (
	"testing"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/accounts/wallet"
)

func TestNewCoinBaseTX(test *testing.T) {
	w := wallet.NewWallet()
	coinBaseTx := NewCoinBaseTX(string(w.GetAddress()), 1.0567)

	if coinBaseTx.Fee != 0 {
		test.Errorf("invalid coin base tx fee:\nactual:\n%f\nexpected:\n0", coinBaseTx.Fee)
	}
	if len(coinBaseTx.VIn) != 1 {
		test.Errorf("invalid coin base tx inputs len:\nactual:\n%d\nexpected:\n1", len(coinBaseTx.VIn))
	}
	if len(coinBaseTx.VIn[0].PreviousTx) != 0 {
		test.Errorf("invalid coin base tx input len of previous tx hash::\nactual:\n%d\nexpected:\n0", len(coinBaseTx.VIn[0].PreviousTx))
	}
	if coinBaseTx.VIn[0].VOut != -1 {
		test.Errorf("invalid coin base tx input out referance:\nactual:\n%d\nexpected:\n-1", coinBaseTx.VIn[0].VOut)
	}
	if coinBaseTx.VIn[0].Signature != nil {
		test.Errorf("invalid coin base tx input signature:\nactual:\n%x\nexpected:\nnil", coinBaseTx.VIn[0].Signature)
	}
	if len(coinBaseTx.VOut) != 1 {
		test.Errorf("invalid coin base tx outs len:\nactual:\n%d\nexpected:\n1", len(coinBaseTx.VOut))
	}
	if coinBaseTx.VOut[0].Value != vars.MINING_REWARD+1.0567 {
		test.Errorf("invalid coin base tx output value:\nactual:\n%f\nexpected:\n%f", coinBaseTx.VOut[0].Value, vars.MINING_REWARD+1.0567)
	}
	expectedHash := coinBaseTx.CalcHash()
	if len(coinBaseTx.Hash) != len(expectedHash) {
		test.Errorf("invalid coin base tx hash:\nactual:\n%x\nexpected:\n%x", coinBaseTx.Hash, expectedHash)
	}
	for i := range coinBaseTx.Hash {
		if coinBaseTx.Hash[i] != expectedHash[i] {
			test.Errorf("invalid coin base tx hash:\nactual:\n%x\nexpected:\n%x", coinBaseTx.Hash, expectedHash)
		}
	}
}
