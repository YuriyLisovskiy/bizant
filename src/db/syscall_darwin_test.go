// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"github.com/stretchr/testify/mock"
)

type mocksyscall struct {
	mock.Mock
}

func (m *mocksyscall) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	args := m.Called(fd, offset, length, prot, flags)
	return args.Get(0).([]byte), args.Error(1)
}
