// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package core

import (
	"log"
	"encoding/hex"

	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	db_pkg "github.com/YuriyLisovskiy/blockchain-go/src/db"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types"
)

type BlockChainIterator struct {
	currentHash []byte
	db          *db_pkg.DB
}

func (i *BlockChainIterator) Next() types.Block {
	encodedBlock, err := i.db.Get(i.currentHash, utils.BLOCKS_BUCKET)
	if err != nil {
		log.Panic(err)
	}
	block := DeserializeBlock(encodedBlock)
	i.currentHash = block.PrevBlockHash
	return block

/*
	err := i.db.View(func(tx *db_pkg.Tx) error {
		b := tx.Bucket(utils.BLOCKS_BUCKET)
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
*/
}

func (i *BlockChainIterator) End() bool {
	return hex.EncodeToString(i.currentHash) == hex.EncodeToString([]byte{})
}
