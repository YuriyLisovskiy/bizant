// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"testing"

	"github.com/boltdb/bolt"
	. "github.com/boltdb/bolt/cmd/bolt"
	"github.com/stretchr/testify/assert"
)

// Ensure that a list of buckets can be retrieved.
func TestBuckets(t *testing.T) {
	SetTestMode(true)
	open(func(db *bolt.DB) {
		db.Update(func(tx *bolt.Tx) error {
			tx.CreateBucket("woojits")
			tx.CreateBucket("widgets")
			tx.CreateBucket("whatchits")
			return nil
		})
		output := run("buckets", db.Path())
		assert.Equal(t, "whatchits\nwidgets\nwoojits", output)
	})
}

// Ensure that an error is reported if the database is not found.
func TestBucketsDBNotFound(t *testing.T) {
	SetTestMode(true)
	output := run("buckets", "no/such/db")
	assert.Equal(t, "stat no/such/db: no such file or directory", output)
}
