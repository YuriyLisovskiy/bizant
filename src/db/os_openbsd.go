// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package db

import (
	"unsafe"
	"syscall"
)

const (
	msAsync      = 1 << iota // perform asynchronous writes
	msSync                   // perform synchronous writes
	msInvalidate             // invalidate cached data
)

func msync(db *DB) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(db.data)), uintptr(db.datasz), msInvalidate)
	if errno != 0 {
		return errno
	}
	return nil
}

func fdatasync(db *DB) error {
	if db.data != nil {
		return msync(db)
	}
	return db.file.Sync()
}
