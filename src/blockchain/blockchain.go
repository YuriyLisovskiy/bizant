package blockchain

import (
	"github.com/boltdb/bolt"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	newBlock := NewBlock(data, lastHash)
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

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockChain(dbFile string) *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(utils.BlocksBucket))
			if err != nil {
				panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	bc := BlockChain{tip, db}
	return &bc
}
