// Copyright (c) 2013 Ben Johnson
// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package arch

// maxMapSize represents the largest mmap size supported by Bolt.
const MaxMapSize = 0xFFFFFFFFFFFF // 256TB

// maxAllocSize is the size used when creating array pointers.
const MaxAllocSize = 0x7FFFFFFF

// Are unaligned load/stores broken on this arch?
var BrokenUnaligned = false
