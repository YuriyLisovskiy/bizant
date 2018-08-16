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
