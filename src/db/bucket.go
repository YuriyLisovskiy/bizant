// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

type Bucket struct {
	*bucket
	name        string
	transaction *Transaction
}

type bucket struct {
	root pgid
}

// Name returns the name of the bucket.
func (b *Bucket) Name() string {
	return b.name
}

// Get retrieves the value for a key in the bucket.
func (b *Bucket) Get(key []byte) []byte {
	return b.Cursor().Get(key)
}

// Cursor creates a new cursor for this bucket.
func (b *Bucket) Cursor() *Cursor {
	return &Cursor{
		transaction: b.transaction,
		root:        b.root,
		stack:       make([]elem, 0),
	}
}
