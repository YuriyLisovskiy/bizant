// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"os"
	"syscall"
)

func fdatasync(f *os.File) error {
	return syscall.Fdatasync(int(f.Fd()))
}
