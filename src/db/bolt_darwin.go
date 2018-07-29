// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"os"
)

var odirect int

// fdatasync flushes written data to a file descriptor.
func fdatasync(f *os.File) error {
	return f.Sync()
}
