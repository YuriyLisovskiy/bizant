// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package consensus

import (
	"fmt"
	"testing"
	"encoding/hex"
)

func TestNewMerkleNode(test *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	// Level 1
	n1 := newMerkleNode(nil, nil, data[0])
	n2 := newMerkleNode(nil, nil, data[1])
	n3 := newMerkleNode(nil, nil, data[2])
	n4 := newMerkleNode(nil, nil, data[2])

	// Level 2
	n5 := newMerkleNode(n1, n2, nil)
	n6 := newMerkleNode(n3, n4, nil)

	// Level 3
	n7 := newMerkleNode(n5, n6, nil)

	if "64b04b718d8b7c5b6fd17f7ec221945c034cfce3be4118da33244966150c4bd4" != hex.EncodeToString(n5.Data) {
		test.Error("Level 1 hash 1 is incorrect")
	}
	if "08bd0d1426f87a78bfc2f0b13eccdf6f5b58dac6b37a7b9441c1a2fab415d76c" != hex.EncodeToString(n6.Data) {
		test.Error("Level 1 hash 2 is incorrect")
	}
	if "4e3e44e55926330ab6c31892f980f8bfd1a6e910ff1ebc3f778211377f35227e" != hex.EncodeToString(n7.Data) {
		test.Error("Root hash is incorrect")
	}
}

func TestNewMerkleTree(test *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	// Level 1
	n1 := newMerkleNode(nil, nil, data[0])
	n2 := newMerkleNode(nil, nil, data[1])
	n3 := newMerkleNode(nil, nil, data[2])
	n4 := newMerkleNode(nil, nil, data[2])

	// Level 2
	n5 := newMerkleNode(n1, n2, nil)
	n6 := newMerkleNode(n3, n4, nil)

	// Level 3
	n7 := newMerkleNode(n5, n6, nil)

	rootHash := fmt.Sprintf("%x", n7.Data)
	mRoot := ComputeMerkleRoot(data)
	if rootHash != fmt.Sprintf("%x", mRoot) {
		test.Error("Merkle tree root hash is incorrect")
	}
}
