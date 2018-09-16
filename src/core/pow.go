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
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/x11"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

type Worker struct {
	block  types.Block
	target *big.Int
}

func NewProofOfWork(block types.Block) Worker {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-vars.TARGET_BITS))
	worker := Worker{block, target}
	return worker
}

func (w *Worker) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			w.block.PrevBlockHash,
			w.block.HashTransactions(),
			utils.IntToHex(w.block.Timestamp),
			utils.IntToHex(int64(vars.TARGET_BITS)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (w *Worker) Run() (int, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < vars.MAX_NONCE {
		if atomic.LoadInt32(&vars.Syncing) == 1 {
			return 0, []byte{}, errors.New("mining interrupt")
		}
		data := w.prepareData(nonce)
		hash = x11.Sum256(data)
		utils.PrintLog(fmt.Sprintf("Mining a new block: %x", hash))
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(w.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n")
	return nonce, hash[:], nil
}

func (w *Worker) Validate() bool {
	var hashInt big.Int
	data := w.prepareData(w.block.Nonce)
	hash := x11.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(w.target) == -1
	return isValid
}
