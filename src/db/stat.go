// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

type Stat struct {
	PageSize          int
	Depth             int
	BranchPageCount   int
	LeafPageCount     int
	OverflowPageCount int
	EntryCount        int
}
