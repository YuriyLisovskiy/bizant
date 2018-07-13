// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package blockchain

import (
	"log"
	"encoding/hex"
	"github.com/boltdb/bolt"
	txPkg "github.com/YuriyLisovskiy/blockchain-go/src/tx"
)

type UTXOSet struct {
	BlockChain BlockChain
}

func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount float64) (float64, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := float64(0)
	db := u.BlockChain.db
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXO_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := txPkg.DeserializeOutputs(v)
			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

func (u UTXOSet) FindUTXO(pubKeyHash []byte) []txPkg.TXOutput {
	var UTXOs []txPkg.TXOutput
	db := u.BlockChain.db
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXO_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := txPkg.DeserializeOutputs(v)
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return UTXOs
}

func (u UTXOSet) CountTransactions() int {
	db := u.BlockChain.db
	counter := 0
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXO_BUCKET))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return counter
}

func (u UTXOSet) Reindex() {
	db := u.BlockChain.db
	bucketName := []byte(UTXO_BUCKET)
	err := db.Batch(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	UTXO := u.BlockChain.FindUTXO()
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

func (u UTXOSet) Update(block Block) {
	db := u.BlockChain.db
	DBMutex.Lock()
	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXO_BUCKET))
		for _, tx := range block.Transactions {
			if tx.IsCoinBase() == false {
				for _, vin := range tx.VIn {
					updatedOuts := txPkg.TXOutputs{}
					outsBytes := b.Get(vin.TxId)
					outs := txPkg.DeserializeOutputs(outsBytes)
					for outIdx, out := range outs.Outputs {
						if outIdx != vin.VOut {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}
					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.TxId)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := b.Put(vin.TxId, updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}
			newOutputs := txPkg.TXOutputs{}
			for _, out := range tx.VOut {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	DBMutex.Unlock()
	if err != nil {
		log.Panic(err)
	}
}
