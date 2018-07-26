// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package core

import (
	"log"
	"encoding/hex"

	"github.com/boltdb/bolt"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
)

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockChainIterator) Next() types.Block {
	var block types.Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utils.BLOCKS_BUCKET))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}

func (i *BlockChainIterator) End() bool {
	return hex.EncodeToString(i.currentHash) == hex.EncodeToString([]byte{})
}
