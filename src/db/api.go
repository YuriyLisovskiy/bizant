// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

//  This file implements simplified api for accessing the database.
//
//  Use this methods if it is not necessary to implement any logic in callback functions,
//  which are used in db.View(...), db.Update(...), and db.Batch(...).

package db

import "errors"

// Get gets data from database bucket by given key.
func (db *DB) Get(key, bucket []byte) ([]byte, error) {
	var result []byte
	err := db.View(func(tx *Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		result = b.Get(key)
		if result == nil {
			return ErrKeyNotFound
		}
		return nil
	})
	return result, err
}

// Put updates bucket in the database with new value by given key.
// If bucket does not exist, Put creates it. Use this method if it is not necessary to write any other logic
// in callback function.
func (db *DB) Put(key, value, bucket []byte, useBatch bool) error {
	function := func(tx *Tx) error {
		var err error
		b := tx.Bucket(bucket)
		if b == nil {
			b, err = tx.CreateBucket(bucket)
			if err != nil {
				return err
			}
		}
		err = b.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	}
	if useBatch {
		return db.Batch(function)
	}
	return db.Update(function)
}

// Same as previous method, but puts an array of values by given keys.
func (db *DB) PutArray(keys, values [][]byte, bucket []byte, useBatch bool) error {
	if len(keys) != len(values) {
		return errors.New("keys len is not equal to values len")
	}
	function := func(tx *Tx) error {
		var err error
		b := tx.Bucket(bucket)
		if b == nil {
			b, err = tx.CreateBucket(bucket)
			if err != nil {
				return err
			}
		}
		for i := range keys {
			err = b.Put(keys[i], values[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
	if useBatch {
		return db.Batch(function)
	}
	return db.Update(function)
}

// Delete removes data from given bucket by key.
func (db *DB) Delete(key, bucket []byte, useBatch bool) error {
	function := func(tx *Tx) error {
		var err error
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		err = b.Delete(key)
		if err != nil {
			return err
		}
		return nil
	}
	if useBatch {
		return db.Batch(function)
	}
	return db.Update(function)
}
