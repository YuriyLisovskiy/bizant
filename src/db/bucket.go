// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

// Bucket represents a collection of key/value pairs inside the database.
// All keys inside the bucket are unique. The Bucket type is not typically used
// directly. Instead the bucket name is typically passed into the Get(), Put(),
// or Delete() functions.
type Bucket struct {
	*bucket
	name        string
	transaction *Transaction
}

// bucket represents the on-file representation of a bucket.
type bucket struct {
	root     pgid
	sequence uint64
}

// Name returns the name of the bucket.
func (b *Bucket) Name() string {
	return b.name
}

// cursor creates a new cursor for this bucket.
func (b *Bucket) cursor() *Cursor {
	return &Cursor{
		transaction: b.transaction,
		root:        b.root,
		stack:       make([]pageElementRef, 0),
	}
}
