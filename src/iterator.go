package src

import "github.com/boltdb/bolt"

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}
	return bci
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block
	err := iterator.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
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
