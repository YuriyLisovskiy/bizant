// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package types

import (
	"log"
	"bytes"
	"encoding/gob"

	"github.com/YuriyLisovskiy/blockchain-go/src/consensus"
)

type Block struct {
	Timestamp     int64
	Transactions  []Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}

func (b Block) HashTransactions() []byte {
	var transactions [][]byte
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	return consensus.ComputeMerkleRoot(transactions)
}

func (b Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}
