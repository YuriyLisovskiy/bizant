// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

import (
	"io"
	"os"
)

type Syncer interface {
	Sync() error
}

type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	Stat() (fi os.FileInfo, err error)
}

type Writer interface {
	io.Writer
	io.Closer
	Syncer
}

type FileLock interface {
	Unlock() error
}

type Env interface {
	OpenReader(name string) (f Reader, err error)
	OpenWriter(name string) (f Writer, err error)
	FileExists(name string) bool
	GetFileSize(name string) (size uint64, err error)
	RenameFile(oldName, newName string) error
	RemoveFile(name string) error
	ListDir(name string) (res []string, err error)
	CreateDir(name string) error
	RemoveDir(name string) error
	LockFile(name string) (fl FileLock, err error)
}

var DefaultEnv = StdEnv{}

type StdEnv struct{}

func (StdEnv) OpenReader(name string) (f Reader, err error) {
	return os.OpenFile(name, os.O_RDONLY, 0)
}

func (StdEnv) OpenWriter(name string) (f Writer, err error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func (StdEnv) FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func (StdEnv) GetFileSize(name string) (size uint64, err error) {
	var fi os.FileInfo
	fi, err = os.Stat(name)
	if err == nil {
		size = uint64(fi.Size())
	}
	return
}

func (StdEnv) RenameFile(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (StdEnv) RemoveFile(name string) error {
	return os.Remove(name)
}

func (StdEnv) ListDir(name string) (res []string, err error) {
	var f *os.File
	f, err = os.Open(name)
	if err != nil {
		return
	}
	res, err = f.Readdirnames(0)
	f.Close()
	return
}

func (StdEnv) CreateDir(name string) error {
	return os.Mkdir(name, 0755)
}

func (StdEnv) RemoveDir(name string) error {
	return os.RemoveAll(name)
}

func (StdEnv) LockFile(name string) (fl FileLock, err error) {
	var f *os.File
	f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	err = setFileLock(f, true)
	if err == nil {
		fl = &stdFileLock{f: f}
	} else {
		f.Close()
	}
	return
}

type stdFileLock struct {
	f *os.File
}

func (fl *stdFileLock) Unlock() (err error) {
	err = setFileLock(fl.f, false)
	fl.f.Close()
	return
}
