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
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/groestl"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/jh"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/keccak"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/luffa"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/shavite"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/simd"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/skein"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/sha3/utils"
)

// Hash contains the state objects
// required to perform the x11.Hash.
type Hash struct {
	tha [64]byte
	thb [64]byte

	blake    utils.Digest
	bmw      utils.Digest
	cubehash utils.Digest
	echo     utils.Digest
	groestl  utils.Digest
	jh       utils.Digest
	keccak   utils.Digest
	luffa    utils.Digest
	shavite  utils.Digest
	simd     utils.Digest
	skein    utils.Digest
}

// New returns a new object to compute a x11 hash.
func New() *Hash {
	ref := &Hash{}

	ref.blake = blake.New()
	ref.bmw = bmw.New()
	ref.cubehash = cubehash.New()
	ref.echo = echo.New()
	ref.groestl = groestl.New()
	ref.jh = jh.New()
	ref.keccak = keccak.New()
	ref.luffa = luffa.New()
	ref.shavite = shavite.New()
	ref.simd = simd.New()
	ref.skein = skein.New()

	return ref
}

// Hash computes the hash from the src bytes and stores the result in dst.
func (ref *Hash) Sum(src []byte) [64]byte {
	ta := ref.tha[:]
	tb := ref.thb[:]

	ref.blake.Write(src)
	ref.blake.Close(tb, 0, 0)

	ref.bmw.Write(tb)
	ref.bmw.Close(ta, 0, 0)

	ref.groestl.Write(ta)
	ref.groestl.Close(tb, 0, 0)

	ref.skein.Write(tb)
	ref.skein.Close(ta, 0, 0)

	ref.jh.Write(ta)
	ref.jh.Close(tb, 0, 0)

	ref.keccak.Write(tb)
	ref.keccak.Close(ta, 0, 0)

	ref.luffa.Write(ta)
	ref.luffa.Close(tb, 0, 0)

	ref.cubehash.Write(tb)
	ref.cubehash.Close(ta, 0, 0)

	ref.shavite.Write(ta)
	ref.shavite.Close(tb, 0, 0)

	ref.simd.Write(tb)
	ref.simd.Close(ta, 0, 0)

	ref.echo.Write(ta)
	ref.echo.Close(tb, 0, 0)

	var res [64]byte
	copy(res[:], tb)

	return res
}
