// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package core

import (
	"testing"

	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
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
