// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

package bolt

import "os"

func fdatasync(f *os.File) error {
	return f.Sync()
}
