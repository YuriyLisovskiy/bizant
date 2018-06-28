package blockchain

import (
	"time"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, transaction := range block.Transactions {
		txHashes = append(txHashes, transaction.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(byteData []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(byteData))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}
