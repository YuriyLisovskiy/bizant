package blockchain

import (
	"github.com/boltdb/bolt"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	txPkg "github.com/YuriyLisovskiy/blockchain-go/src/tx"
	"fmt"
	"os"
	"log"
	"encoding/hex"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	newBlock := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
}

func (bc *BlockChain) CloseDB() {
	defer bc.db.Close()
}

func NewGenesisBlock(coinBase *Transaction) *Block {
	return NewBlock([]*Transaction{coinBase}, []byte{})
}

func NewBlockChain(address string) *BlockChain {
	if !utils.DBExists(utils.DBFile) {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	bc := BlockChain{tip, db}
	return &bc
}

func CreateBlockChain(address string) *BlockChain {
	if utils.DBExists(utils.DBFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(utils.DBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		genesis := NewGenesisBlock(NewCoinBaseTX(address, utils.GenesisCoinbaseData))
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
	bc := BlockChain{tip, db}
	return &bc
}

func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.VOut {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.VIn {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.TxId)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VOut)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *BlockChain) FindUTXO(address string) []txPkg.TXOutput {
	var UTXOs []txPkg.TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.VOut {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.VOut {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}
