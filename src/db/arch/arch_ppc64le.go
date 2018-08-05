// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

// +build ppc64le

package arch

// maxMapSize represents the largest mmap size supported by Bolt.
const MaxMapSize = 0xFFFFFFFFFFFF // 256TB

// maxAllocSize is the size used when creating array pointers.
const MaxAllocSize = 0x7FFFFFFF

// Are unaligned load/stores broken on this arch?
var BrokenUnaligned = false
