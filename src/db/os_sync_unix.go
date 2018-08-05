// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

// +build !windows,!plan9,!linux,!openbsd

package db

// fdatasync flushes written data to a file descriptor.
func fdatasync(db *DB) error {
	return db.file.Sync()
}
