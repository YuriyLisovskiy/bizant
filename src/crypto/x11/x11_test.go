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
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSum512(t *testing.T) {
	out := [64]byte{}

	for i := range Sum512_Data {
		ln := len(Sum512_Data[i].out)
		dest := make([]byte, ln)

		out = Sum512(Sum512_Data[i].in[:])
		if ln != hex.Encode(dest, out[:]) {
			t.Errorf("%s: invalid length", Sum512_Data[i])
		}
		if !bytes.Equal(dest[:], Sum512_Data[i].out[:]) {
			t.Errorf("%s: invalid hash", Sum512_Data[i].id)
		}
	}
}

var Sum512_Data = []struct {
	id  string
	in  []byte
	out []byte
}{
	{
		"Empty",
		[]byte(""),
		[]byte("51b572209083576ea221c27e62b4e22063257571ccb6cc3dc3cd17eb67584eba3dfd9d129b61e0d802866f5d09ab2c280ca07242380a811d10bb0437ce546065"),
	},
	{
		"Dash",
		[]byte("DASH"),
		[]byte("fe809ebca8753d907f6ad32cdcf8e5c4e090d7bece5df35b2147e10b88c12d26578b18d97bd9ca71c35549cd04fc3449a7c910814808133a2f976c42fc28f2df"),
	},
	{
		"Fox",
		[]byte("The quick brown fox jumps over the lazy dog"),
		[]byte("534536a4e4f16b32447f02f77200449dc2f23b532e3d9878fe111c9de666bc5cafc61ae1a2884127d00d897065528dc35d2ea9222d95e8f6e94e1f0b52bdcddc"),
	},
}

func TestSum256(t *testing.T) {
	out := [32]byte{}

	for i := range Sum256_Data {
		ln := len(Sum256_Data[i].out)
		dest := make([]byte, ln)

		out = Sum256(Sum256_Data[i].in[:])
		if ln != hex.Encode(dest, out[:]) {
			t.Errorf("%s: invalid length", Sum256_Data[i])
		}
		if !bytes.Equal(dest[:], Sum256_Data[i].out[:]) {
			t.Errorf("%s: invalid hash", Sum256_Data[i].id)
		}
	}
}

var Sum256_Data = []struct {
	id  string
	in  []byte
	out []byte
}{
	{
		"The great experiment continues",
		[]byte("The great experiment continues"),
		[]byte("e05103283876cfa7254683f678f0b1a4c3621ffdd51b78bad2fa134b4875c936"),
	},
	{
		"Empty",
		[]byte(""),
		[]byte("51b572209083576ea221c27e62b4e22063257571ccb6cc3dc3cd17eb67584eba"),
	},
	{
		"Dash",
		[]byte("DASH"),
		[]byte("fe809ebca8753d907f6ad32cdcf8e5c4e090d7bece5df35b2147e10b88c12d26"),
	},
	{
		"Fox",
		[]byte("The quick brown fox jumps over the lazy dog"),
		[]byte("534536a4e4f16b32447f02f77200449dc2f23b532e3d9878fe111c9de666bc5c"),
	},
}
