package blockchain

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"errors"
	"encoding/hex"
	"crypto/ecdsa"
	"github.com/boltdb/bolt"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	txPkg "github.com/YuriyLisovskiy/blockchain-go/src/tx"
//	"encoding/json"
	"encoding/json"
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
	cbTx := NewCoinBaseTX(address, 0, utils.GenesisCoinbaseData)
	genesis := NewGenesisBlock(cbTx)
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(utils.BlocksBucket))
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
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return BlockChain{tip, db}
}

func (bc *BlockChain) AddBlock(block Block) {
	DBMutex.Lock()
	err := bc.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
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
	DBMutex.Unlock()
	if err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) GetBestHeight() int {
	var lastBlock Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
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

func (bc *BlockChain) GetBlock(blockHash []byte) (Block, error) {
	var block Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
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

func (bc *BlockChain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		blocks = append(blocks, block.Hash)
	}
	return blocks
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for !bci.End() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return tx, nil
			}
		}
	}
	return Transaction{}, errors.New("transaction is not found")
}

func (bc *BlockChain) FindUTXO() map[string]txPkg.TXOutputs {
	UTXO := make(map[string]txPkg.TXOutputs)
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
	DBMutex.Lock()
	err := bc.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		bc.tip = b.Get([]byte("l"))
		return nil
	})
	DBMutex.Unlock()
	if err != nil {
		log.Panic(err)
	}

	return BlockChainIterator{bc.tip, bc.db}
}

func (bc *BlockChain) MineBlock(minerAddress string, transactions []Transaction) (Block, error) {
	var lastHash []byte
	var lastHeight int
	fees := 0.0
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {

			// TODO: send an error to transaction's author

			utils.PrintLog(fmt.Sprintf("Invalid transaction %x\n", tx.ID))

			data, err := json.MarshalIndent(tx, "", "  ")
			if err == nil {
				fmt.Println(string(data))
			}

		} else {
			fees += tx.Fee
		}
	}
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
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
		return Block{}, err
	}
	DBMutex.Lock()
	err = bc.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
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
	DBMutex.Unlock()
	if err != nil {
		log.Panic(err)
	}
	return bc.GetBlock(newBlock.Hash)
}

func (bc *BlockChain) SignTransaction(tx Transaction, privKey ecdsa.PrivateKey) Transaction {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.VIn {
		prevTX, err := bc.FindTransaction(vin.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.VIn {
		prevTX, err := bc.FindTransaction(vin.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}

func (bc *BlockChain) CloseDB(Defer bool) {
	if Defer {
		defer bc.db.Close()
	}
	bc.db.Close()
}
