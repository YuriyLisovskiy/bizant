// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"syscall"
)

type _syscall interface {
	Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
	Munmap([]byte) error
}

type syssyscall struct{}

func (o *syssyscall) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return syscall.Mmap(fd, offset, length, prot, flags)
}

func (o *syssyscall) Munmap(b []byte) error {
	return syscall.Munmap(b)
}
