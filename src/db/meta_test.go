// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensure that meta with bad magic is invalid.
func TestMetaValidateMagic(t *testing.T) {
	m := &meta{magic: 0x01234567}
	assert.Equal(t, m.validate(), ErrInvalid)
}

// Ensure that meta with a bad version is invalid.
func TestMetaValidateVersion(t *testing.T) {
	m := &meta{magic: magic, version: 200}
	assert.Equal(t, m.validate(), ErrVersionMismatch)
}
