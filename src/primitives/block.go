// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package primitives

import (
	"log"
	"time"
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

func NewBlock(transactions []Transaction, prevBlockHash []byte, height int) (Block, error) {
	block := Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash, err := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block, err
}

func NewGenesisBlock(coinBase Transaction) Block {
	block, _ := NewBlock([]Transaction{coinBase}, []byte{}, 0)
	return block
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

func DeserializeBlock(d []byte) Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return block
}
