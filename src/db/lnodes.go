// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

type lnodes []lnode

// replace replaces the node at the given index with a new key/value size.
func (s lnodes) replace(key, value []byte, index int) lnodes {
	n := &s[index]
	n.pos = 0
	n.ksize = len(key)
	n.vsize = len(value)
	return s
}

// insert places a new node at the given index with a key/value size.
func (s lnodes) insert(key, value []byte, index int) lnodes {
	return append(s[0:index], lnode{ksize: len(key), vsize: len(value)}, s[index:len(s)])
}
