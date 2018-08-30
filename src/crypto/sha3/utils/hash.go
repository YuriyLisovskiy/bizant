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

package utils

import "io"

type Hash interface {
	// Write (via the embedded io.Writer interface) adds more
	// data to the running hash. It never returns an error.
	io.Writer

	// Reset resets the Hash to its initial state.
	Reset()

	// Sum appends the current hash to dst and returns the result
	// as a slice. It does not change the underlying hash state.
	Sum(dst []byte) []byte

	// Size returns the number of bytes Sum will return.
	Size() int

	// BlockSize returns the hash's underlying block size.
	// The Write method must be able to accept any amount
	// of data, but it may operate more efficiently if
	// all writes are a multiple of the block size.
	BlockSize() int
}
