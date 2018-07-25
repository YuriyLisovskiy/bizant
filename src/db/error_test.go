// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensure that nested errors are appropriately formatted.
func TestError(t *testing.T) {
	e := &Error{"one error", &Error{"two error", nil}}
	assert.Equal(t, e.Error(), "one error: two error")
}
