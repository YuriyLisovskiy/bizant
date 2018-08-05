// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package core

import (
	"log"
	"time"
	"bytes"
	"encoding/gob"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types/tx_io"
)

func NewBlock(transactions []types.Transaction, prevBlockHash []byte, height int) (types.Block, error) {
	block := types.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Height:        height,
	}
	pow := NewProofOfWork(block)
	nonce, hash, err := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block, err
}

func NewGenesisBlock(coinBase types.Transaction) (types.Block, error) {
	return NewBlock([]types.Transaction{coinBase}, []byte{}, 0)
}

func DeserializeBlock(d []byte) types.Block {
	var block types.Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return block
}

func NewCoinBaseTX(to string, fees float64) types.Transaction {
	txIn := tx_io.TXInput{PreviousTx: []byte{}, VOut: -1, Signature: nil}
	txOut := tx_io.NewTXOutput(vars.MINING_REWARD+fees, to)
	tx := types.Transaction{
		Hash:        nil,
		VIn:       []tx_io.TXInput{txIn},
		VOut:      []tx_io.TXOutput{txOut},
		Timestamp: time.Now().Unix(),
		Fee:       0,
	}
	tx.Hash = tx.CalcHash()
	return tx
}

func DeserializeTransaction(data []byte) types.Transaction {
	var transaction types.Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}
