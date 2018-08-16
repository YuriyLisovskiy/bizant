// Copyright (c) 2013 Ben Johnson
// Licensed under the MIT software license.
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

package arch

import "unsafe"

// maxMapSize represents the largest mmap size supported by Bolt.
const MaxMapSize = 0x7FFFFFFF // 2GB

// maxAllocSize is the size used when creating array pointers.
const MaxAllocSize = 0xFFFFFFF

// Are unaligned load/stores broken on this arch?
var BrokenUnaligned bool

func init() {
	// Simple check to see whether this arch handles unaligned load/stores
	// correctly.

	// ARM9 and older devices require load/stores to be from/to aligned
	// addresses. If not, the lower 2 bits are cleared and that address is
	// read in a jumbled up order.

	// See http://infocenter.arm.com/help/index.jsp?topic=/com.arm.doc.faqs/ka15414.html

	raw := [6]byte{0xfe, 0xef, 0x11, 0x22, 0x22, 0x11}
	val := *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&raw)) + 2))
	brokenUnaligned = val != 0x11222211
}
