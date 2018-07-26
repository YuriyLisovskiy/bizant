// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"os"
)

// Open creates and opens a database at the given path.
// If the file does not exist then it will be created automatically.
func Open(path string, mode os.FileMode) (*DB, error) {
	db := &DB{}
	if err := db.Open(path, mode); err != nil {
		return nil, err
	}
	return db, nil
}
