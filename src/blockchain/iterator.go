package blockchain

import (
	"github.com/boltdb/bolt"
	"github.com/YuriyLisovskiy/BlockChainGo/src/utils"
)

type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *BlockChain) Iterator() *Iterator {
	bci := &Iterator{bc.tip, bc.db}
	return bci
}

func (iterator *Iterator) Next() *Block {
	var block *Block
	err := iterator.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BlocksBucket))
		encodedBlock := b.Get(iterator.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		return nil
	}
	iterator.currentHash = block.PrevBlockHash
	return block
}
