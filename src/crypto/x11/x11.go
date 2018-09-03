// ISC License
// 
// Copyright (c) 2017-2018 The DashX developers
// Copyright (c) 2016 The Dash developers
// 
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
// 
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//
//
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

package x11

import (
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/blake"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/bmw"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/cubehash"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/echo"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/groestl512"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/jh"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/keccak"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/luffa"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/shavite"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/simd"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/skein"
)

var (
	tha [64]byte
	thb [64]byte

	blakeHash = blake.New()
	bmwHash = bmw.New()
	groestlHash = groestl512.New()
	skeinHash = skein.New()
	jhHash = jh.New()
	keccakHash = keccak.New()
	luffaHash = luffa.New()
	cubehashHash = cubehash.New()
	shaviteHash = shavite.New()
	simdHash = simd.New()
	echoHash = echo.New()
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
