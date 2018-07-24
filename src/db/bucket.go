// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

var (
	InvalidBucketError = &Error{"invalid bucket", nil}
)

type bucketid uint32

type Bucket struct {
	*bucket
	name        string
	transaction *Transaction
	cursors     []*Cursor
}

type bucket struct {
	id    bucketid
	flags uint32
	root  pgid
}

// Get retrieves the value for a key in the bucket.
func (b *Bucket) Get(key []byte) []byte {
	return b.cursor().Get(key)
}

// Cursor creates a new cursor for this bucket.
func (b *Bucket) Cursor() *Cursor {
	c := b.cursor()
	b.cursors = append(b.cursors, c)
	return c
}

// cursor creates a new untracked cursor for this bucket.
func (b *Bucket) cursor() *Cursor {
	return &Cursor{
		bucket: b,
		stack:  make([]elem, 0),
	}
}
