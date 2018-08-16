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
