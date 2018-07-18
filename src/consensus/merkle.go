// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package consensus

import "crypto/sha256"

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func ComputeMerkleRoot(data [][]byte) []byte {
	var nodes []MerkleNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, datum := range data {
		node := newMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := newMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}
	return nodes[0].Data
}

func newMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}
	mNode.Left = left
	mNode.Right = right
	return &mNode
}

