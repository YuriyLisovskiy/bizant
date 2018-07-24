// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"os"
)

type _os interface {
	OpenFile(name string, flag int, perm os.FileMode) (file file, err error)
	Stat(name string) (fi os.FileInfo, err error)
	Getpagesize() int
}

type sysos struct{}

func (o *sysos) OpenFile(name string, flag int, perm os.FileMode) (file file, err error) {
	return os.OpenFile(name, flag, perm)
}

func (o *sysos) Stat(name string) (fi os.FileInfo, err error) {
	return os.Stat(name)
}

func (o *sysos) Getpagesize() int {
	return os.Getpagesize()
}
