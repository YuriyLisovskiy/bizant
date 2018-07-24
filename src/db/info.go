// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

// Info contains information about the database.
type Info struct {
	MapSize           int
	LastPageID        int
	LastTransactionID int
	MaxReaders        int
	ReaderCount       int
}
