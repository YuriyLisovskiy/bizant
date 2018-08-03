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

	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	db_pkg "github.com/YuriyLisovskiy/blockchain-go/src/db"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types/tx_io"
)

type BlockChain struct {
	tip []byte
	db  *db_pkg.DB
}

func CreateBlockChain(address, nodeID string) BlockChain {
	utils.DBFile = fmt.Sprintf(utils.DBFile, nodeID)
	if utils.DBExists(utils.DBFile) {
		fmt.Printf("%s already exists.\n", utils.DBFile)
		os.Exit(1)
	}
	cbTx := NewCoinBaseTX(address, 0)
	genesis, err := NewGenesisBlock(cbTx)
	if err != nil {
		log.Panic(err)
	}
	db, err := db_pkg.Open(utils.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	keys := [][]byte{
		genesis.Hash,
		utils.LAST_BLOCK_HASH,
	}
	values := [][]byte{
		genesis.Serialize(),    // genesis.Hash
		genesis.Hash,           // utils.LAST_BLOCK_HASH
	}
	err = db.PutArray(keys, values, utils.BLOCKS_BUCKET, false)

/*
	err = db.Update(func(tx *db_pkg.Tx) error {
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
		return nil
	})
*/

	if err != nil {
		log.Panic(err)
	}
	return BlockChain{genesis.Hash, db}
}

func NewBlockChain(nodeID string) BlockChain {
	utils.DBFile = fmt.Sprintf(utils.DBFile, nodeID)
	if utils.DBExists(utils.DBFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	db, err := db_pkg.Open(utils.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	tip, err := db.Get(utils.LAST_BLOCK_HASH, utils.BLOCKS_BUCKET)

//	err = db.View(func(tx *db_pkg.Tx) error {
//		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
//		tip = b.Get([]byte("l"))
//		return nil
//	})
	if err != nil {
		log.Panic(err)
	}
	return BlockChain{tip, db}
}

// AddBlock writes given block to the database if it does not exist.
func (bc *BlockChain) AddBlock(block types.Block) {

	// Check if given block already exists in the database.
	// If exists then returns from function, else performs adding new block logic.
	blockInDb, err := bc.db.Get(block.Hash, utils.BLOCKS_BUCKET)
	if blockInDb != nil {
		return
	}
	if err != nil && err != db_pkg.ErrKeyNotFound {
		log.Panic(err)
	}

	// Lock thread while changing database content.
	vars.DBMutex.Lock()
	err = bc.db.Batch(func(tx *db_pkg.Tx) error {
		b := tx.Bucket(utils.BLOCKS_BUCKET)

		// Write new block to the database
		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		// Check if given block is the newest one.
		lastHash := b.Get(utils.LAST_BLOCK_HASH)
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)
		if block.Height > lastBlock.Height {
			err = b.Put(utils.LAST_BLOCK_HASH, block.Hash)
			if err != nil {
				log.Panic(err)
			}
			//	bc.tip = block.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	vars.DBMutex.Unlock()
}

// GetBestHeight returns the height of the last block.
func (bc *BlockChain) GetBestHeight() int {
	var lastBlock types.Block
	err := bc.db.View(func(tx *db_pkg.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))

		// Retrieve the link to the last block is written in the database.
		lastHash := b.Get(utils.LAST_BLOCK_HASH)
		if lastHash == nil {
			return errors.New("bc.GetBestHeight: last block hash does not exist")
		}

		// Get the last block from the database to retrieve its height.
		blockData := b.Get(lastHash)
		if blockData == nil {
			return errors.New("bc.GetBestHeight: last block does not exist or last hash is invalid")
		}
		lastBlock = DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return lastBlock.Height
}

// GetBlock retrieves a block by given hash and deserialize it.
func (bc *BlockChain) GetBlock(blockHash []byte) (types.Block, error) {
	var block types.Block
	blockData, err := bc.db.Get(blockHash, utils.BLOCKS_BUCKET)
	if err != nil {
		return block, err
	}
	return DeserializeBlock(blockData), nil

/*
	err := bc.db.View(func(tx *db_pkg.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		blockData := b.Get(blockHash)
		if blockData == nil {
			return errors.New("block is not found")
		}
		block = DeserializeBlock(blockData)
		return nil
	})
*/
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
			txID := hex.EncodeToString(tx.Hash)
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
					inTxID := hex.EncodeToString(in.PreviousTx)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VOut)
				}
			}
		}
	}
	return UTXO
}

// Iterator creates and returns a new blockchain iterator
func (bc *BlockChain) Iterator() BlockChainIterator {

	// Retrieve last block hash.
	tip, err := bc.db.Get(utils.LAST_BLOCK_HASH, utils.BLOCKS_BUCKET)
	if err != nil {
		log.Panic(err)
	}
	bc.tip = tip
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
		prevTx, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			inputs = append(inputs, tx_io.TXInput{PreviousTx: prevTx, VOut: out, Signature: nil, PubKey: targetWallet.PublicKey})
		}
	}
	from := fmt.Sprintf("%x", targetWallet.GetAddress())
	outputs = append(outputs, tx_io.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, tx_io.NewTXOutput(acc-amount, from)) // a change
	}
	tx := types.Transaction{
		Hash:      nil,
		VIn:       inputs,
		VOut:      outputs,
		Timestamp: time.Now().Unix(),
		Fee:       0,
	}
	tx.Hash = tx.CalcHash()
	tx.Fee = tx.CalculateFee(fee)
	return utxoSet.BlockChain.SignTransaction(tx, targetWallet.PrivateKey)
}

// MineBlock generates new block.
func (bc *BlockChain) MineBlock(minerAddress string, transactions []types.Transaction) (types.Block, error) {
	var lastHash []byte
	var lastHeight int
	fees := 0.0

	// Verify all given transactions
	// If transaction is invalid, ignore it and send an error to its owner
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {

			// TODO: send an error to transaction's author

		} else {
			fees += tx.Fee
		}
	}

	// Retrieve last block height.
	err := bc.db.View(func(tx *db_pkg.Tx) error {
		b := tx.Bucket(utils.BLOCKS_BUCKET)
		if b == nil {
			return errors.New(fmt.Sprintf("bucket '%x' does not exist", utils.BLOCKS_BUCKET))
		}

		// Get a link to the last block.
		lastHash = b.Get([]byte("l"))

		// Get and deserialize the last block.
		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)
		lastHeight = block.Height
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// Add coin base transaction to a list of transaction.
	transactions = append(transactions, NewCoinBaseTX(minerAddress, fees))

	// Generate new block.
	newBlock, err := NewBlock(transactions, lastHash, lastHeight+1)
	if err != nil {
		fmt.Println(err.Error())
		return types.Block{}, err
	}

	// Lock thread for safe database update.
	vars.DBMutex.Lock()
	err = bc.db.Batch(func(tx *db_pkg.Tx) error {
		b := tx.Bucket(utils.BLOCKS_BUCKET)
		if b == nil {
			return errors.New(fmt.Sprintf("bucket '%x' does not exist", utils.BLOCKS_BUCKET))
		}
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(utils.LAST_BLOCK_HASH, newBlock.Hash)
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
			if bytes.Compare(tx.Hash, ID) == 0 {
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
		prevTX, err := bc.FindTransaction(vin.PreviousTx)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Hash)] = prevTX
	}
	return tx.Verify(prevTXs)
}

func (bc *BlockChain) SignTransaction(tx types.Transaction, privKey []byte) types.Transaction {
	prevTXs := make(map[string]types.Transaction)
	for _, vin := range tx.VIn {
		prevTX, err := bc.FindTransaction(vin.PreviousTx)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Hash)] = prevTX
	}
	return tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) CloseDB(Defer bool) {
	if Defer {
		defer bc.db.Close()
	}
	bc.db.Close()
}
