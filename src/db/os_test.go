// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type mockos struct {
	mock.Mock
}

func (m *mockos) OpenFile(name string, flag int, perm os.FileMode) (file file, err error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(*mockfile), args.Error(1)
}

func (m *mockos) Stat(name string) (fi os.FileInfo, err error) {
	args := m.Called(name)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *mockos) Getpagesize() int {
	args := m.Called()
	return args.Int(0)
}
