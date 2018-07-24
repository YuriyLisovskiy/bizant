// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"unsafe"
)

// leafNode represents a node on a leaf page.
type leafNode struct {
	flags    uint16
	keySize  uint16
	dataSize uint32
	data     uintptr // Pointer to the beginning of the data.
}

// key returns a byte slice that of the key data.
func (n *leafNode) key() []byte {
	return (*[MaxKeySize]byte)(unsafe.Pointer(&n.data))[:n.keySize]
}
