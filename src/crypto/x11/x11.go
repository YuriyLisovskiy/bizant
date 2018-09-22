// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package x11 implements X11 hash algorithm which uses eleven different
// hashing functions, chained together, to calculate the block header.
// The eleven algorithms are blake, bmw, groestl, jh, keccak, skein, luffa,
// cubehash, shavite, simd and echo. Originally developed by Dash creator,
// Evan Duffield, in 2014.
package x11

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/blake512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/bmw512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/cubehash512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/echo512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/groestl512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/jh512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/keccak512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/luffa512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/shavite512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/simd512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/skein512"
)

var (
	tha [64]byte
	thb [64]byte

	blakeHash = blake512.New()
	bmwHash = bmw512.New()
	groestlHash = groestl512.New()
	skeinHash = skein512.New()
	jhHash = jh512.New()
	keccakHash = keccak512.New()
	luffaHash = luffa512.New()
	cubehashHash = cubehash512.New()
	shaviteHash = shavite512.New()
	simdHash = simd512.New()
	echoHash = echo512.New()
)

// Hash computes the hash from the src bytes and returns 32-byte hash.
func Sum256(src []byte) [32]byte {
	ta := tha[:]
	tb := thb[:]

	blakeHash.Write(src)
	blakeHash.Close(tb, 0, 0)

	bmwHash.Write(tb)
	bmwHash.Close(ta, 0, 0)

	groestlHash.Write(ta)
	groestlHash.Close(tb, 0, 0)

	skeinHash.Write(tb)
	skeinHash.Close(ta, 0, 0)

	jhHash.Write(ta)
	jhHash.Close(tb, 0, 0)

	keccakHash.Write(tb)
	keccakHash.Close(ta, 0, 0)

	luffaHash.Write(ta)
	luffaHash.Close(tb, 0, 0)

	cubehashHash.Write(tb)
	cubehashHash.Close(ta, 0, 0)

	shaviteHash.Write(ta)
	shaviteHash.Close(tb, 0, 0)

	simdHash.Write(tb)
	simdHash.Close(ta, 0, 0)

	echoHash.Write(ta)
	echoHash.Close(tb, 0, 0)

	var res [32]byte
	copy(res[:], tb)

	return res
}

// Hash computes the hash from the src bytes and returns 64-byte hash.
func Sum512(src []byte) [64]byte {
	ta := tha[:]
	tb := thb[:]

	blakeHash.Write(src)
	blakeHash.Close(tb, 0, 0)

	bmwHash.Write(tb)
	bmwHash.Close(ta, 0, 0)

	groestlHash.Write(ta)
	groestlHash.Close(tb, 0, 0)

	skeinHash.Write(tb)
	skeinHash.Close(ta, 0, 0)

	jhHash.Write(ta)
	jhHash.Close(tb, 0, 0)

	keccakHash.Write(tb)
	keccakHash.Close(ta, 0, 0)

	luffaHash.Write(ta)
	luffaHash.Close(tb, 0, 0)

	cubehashHash.Write(tb)
	cubehashHash.Close(ta, 0, 0)

	shaviteHash.Write(ta)
	shaviteHash.Close(tb, 0, 0)

	simdHash.Write(tb)
	simdHash.Close(ta, 0, 0)

	echoHash.Write(ta)
	echoHash.Close(tb, 0, 0)

	var res [64]byte
	copy(res[:], tb)

	return res
}
