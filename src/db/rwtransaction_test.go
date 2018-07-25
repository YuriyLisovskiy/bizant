// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure that a RWTransaction can be retrieved.
func TestRWTransaction(t *testing.T) {
	withOpenDB(func(db *DB, path string) {
		txn, err := db.RWTransaction()
		assert.NotNil(t, txn)
		assert.NoError(t, err)
	})
}

// Ensure that a bucket can be created and retrieved.
func TestTransactionCreateBucket(t *testing.T) {
	withOpenDB(func(db *DB, path string) {
		// Create a bucket.
		err := db.CreateBucket("widgets")
		assert.NoError(t, err)

		// Read the bucket through a separate transaction.
		b, err := db.Bucket("widgets")
		assert.NotNil(t, b)
		assert.NoError(t, err)
	})
}
