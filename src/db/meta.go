// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

var (
	InvalidError         = &Error{"Invalid database", nil}
	VersionMismatchError = &Error{"version mismatch", nil}
)

const magic uint32 = 0xDEADC0DE
const version uint32 = 1

type meta struct {
	magic    uint32
	version  uint32
	pageSize uint32
	pgid     pgid
	free     pgid
	txnid    txnid
	sys      bucket
}

// validate checks the marker bytes and version of the meta page to ensure it matches this binary.
func (m *meta) validate() error {
	if m.magic != magic {
		return InvalidError
	} else if m.version != Version {
		return VersionMismatchError
	}
	return nil
}

// copy copies one meta object to another.
func (m *meta) copy(dest *meta) {
	dest.pageSize = m.pageSize
	dest.pgid = m.pgid
	dest.txnid = m.txnid
	dest.sys = m.sys
}
