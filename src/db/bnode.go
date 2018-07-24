// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"unsafe"
)

const bnodeSize = int(unsafe.Sizeof(lnode{}))

// bnode represents a node on a branch page.
type bnode struct {
	pos   uint32
	ksize uint32
	pgid  pgid
}

// key returns a byte slice of the node key.
func (n *bnode) key() []byte {
	return (*[MaxKeySize]byte)(unsafe.Pointer(&n))[n.pos : n.pos+n.ksize]
}
