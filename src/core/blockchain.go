// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package core

import (
	"os"
	"fmt"
	"log"
	"time"
	"bytes"
	"errors"
	"encoding/hex"
	"crypto/ecdsa"
	//	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types/tx_io"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func CreateBlockChain(address, nodeID string) BlockChain {
	utils.DBFile = fmt.Sprintf(utils.DBFile, nodeID)
	if utils.DBExists(utils.DBFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	cbTx := NewCoinBaseTX(address, 0, utils.GENESIS_COINBASE_DATA)
	genesis, err := NewGenesisBlock(cbTx)
	if err != nil {
		log.Panic(err)
	}
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(utils.BLOCKS_BUCKET))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return BlockChain{tip, db}
}

func NewBlockChain(nodeID string) BlockChain {
	utils.DBFile = fmt.Sprintf(utils.DBFile, nodeID)
	if utils.DBExists(utils.DBFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return BlockChain{tip, db}
}

func (bc *BlockChain) AddBlock(block types.Block) {
	vars.DBMutex.Lock()
	err := bc.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		blockInDb := b.Get(block.Hash)
		if blockInDb != nil {
			return nil
		}
		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}
		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)
		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			//	bc.tip = block.Hash
		}
		return nil
	})
	vars.DBMutex.Unlock()
	if err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) GetBestHeight() int {
	var lastBlock types.Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return lastBlock.Height
}

func (bc *BlockChain) GetBlock(blockHash []byte) (types.Block, error) {
	var block types.Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		blockData := b.Get(blockHash)
		if blockData == nil {
			return errors.New("block is not found")
		}
		block = DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		return block, err
	}
	return block, nil
}

func (bc *BlockChain) GetBlockHashes(height int) [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		if block.Height <= height {
			break
		}
		blocks = append(blocks, block.Hash)
	}
	return blocks
}

func (bc *BlockChain) FindUTXO() map[string]tx_io.TXOutputs {
	UTXO := make(map[string]tx_io.TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.VOut {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.VIn {
					inTxID := hex.EncodeToString(in.TxId)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VOut)
				}
			}
		}
	}
	return UTXO
}

func (bc *BlockChain) Iterator() BlockChainIterator {
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		bc.tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return BlockChainIterator{bc.tip, bc.db}
}

func NewUTXOTransaction(targetWallet *wallet.Wallet, to string, amount, fee float64, utxoSet *UTXOSet) types.Transaction {
	var inputs []tx_io.TXInput
	var outputs []tx_io.TXOutput
	pubKeyHash := wallet.HashPubKey(targetWallet.PublicKey)
	acc, validOutputs := utxoSet.FindSpendableOutputs(pubKeyHash, amount, )
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			inputs = append(inputs, tx_io.TXInput{TxId: txID, VOut: out, Signature: nil, PubKey: targetWallet.PublicKey})
		}
	}
	from := fmt.Sprintf("%x", targetWallet.GetAddress())
	outputs = append(outputs, *tx_io.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *tx_io.NewTXOutput(acc-amount, from)) // a change
	}
	tx := types.Transaction{
		ID:        nil,
		VIn:       inputs,
		VOut:      outputs,
		Timestamp: time.Now().Unix(),
		Fee:       0,
	}
	tx.ID = tx.Hash()
	tx.Fee = tx.CalculateFee(fee)
	return utxoSet.BlockChain.SignTransaction(tx, targetWallet.PrivateKey)
}

func (bc *BlockChain) MineBlock(minerAddress string, transactions []types.Transaction) (types.Block, error) {
	var lastHash []byte
	var lastHeight int
	fees := 0.0
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {

			// TODO: send an error to transaction's author

		} else {
			fees += tx.Fee
		}
	}
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		lastHash = b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)
		lastHeight = block.Height
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	transactions = append(transactions, NewCoinBaseTX(minerAddress, fees, ""))
	newBlock, err := NewBlock(transactions, lastHash, lastHeight+1)
	if err != nil {
		fmt.Println(err.Error())
		return types.Block{}, err
	}
	vars.DBMutex.Lock()
	err = bc.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	vars.DBMutex.Unlock()
	if err != nil {
		log.Panic(err)
	}
	return bc.GetBlock(newBlock.Hash)
}

func (bc *BlockChain) FindTransaction(ID []byte) (types.Transaction, error) {
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return tx, nil
			}
		}
	}
	return types.Transaction{}, errors.New("transaction is not found")
}

func (bc *BlockChain) VerifyTransaction(tx types.Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	prevTXs := make(map[string]types.Transaction)
	for _, vin := range tx.VIn {
		prevTX, err := bc.FindTransaction(vin.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}

func (bc *BlockChain) SignTransaction(tx types.Transaction, privKey ecdsa.PrivateKey) types.Transaction {
	prevTXs := make(map[string]types.Transaction)
	for _, vin := range tx.VIn {
		prevTX, err := bc.FindTransaction(vin.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) CloseDB(Defer bool) {
	if Defer {
		defer bc.db.Close()
	}
	bc.db.Close()
}
