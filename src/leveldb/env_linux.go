// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import (
	"os"
	"syscall"
	"unsafe"
)

func setFileLock(f *os.File, lock bool) (err error) {
	lType := syscall.F_UNLCK
	if lock {
		lType = syscall.F_WRLCK
	}
	k := struct {
		Type   uint32
		Whence uint32
		Start  uint64
		Len    uint64
		Pid    uint32
	}{
		Type:   uint32(lType),
		Whence: uint32(os.SEEK_SET),
		Start:  0,
		Len:    0, // lock the entire file.
		Pid:    uint32(os.Getpid()),
	}
	_, _, err = syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(syscall.F_SETLK), uintptr(unsafe.Pointer(&k)))
	return
}
