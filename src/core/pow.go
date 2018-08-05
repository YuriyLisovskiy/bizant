// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package core

import (
	"fmt"
	"bytes"
	"errors"
	"math/big"
	"sync/atomic"
	"crypto/sha256"

	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
)

type ProofOfWork struct {
	block  types.Block
	target *big.Int
}

func NewProofOfWork(block types.Block) ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-vars.TARGET_BITS))
	pow := ProofOfWork{block, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(vars.TARGET_BITS)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < vars.MAX_NONCE {
		if atomic.LoadInt32(&vars.Syncing) == 1 {
			return 0, []byte{}, errors.New("mining interrupt")
		}
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		utils.PrintLog(fmt.Sprintf("Mining a new block: %x", hash))
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n")
	return nonce, hash[:], nil
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
