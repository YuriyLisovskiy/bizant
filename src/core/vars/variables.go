// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package vars

import "sync"

var (
	Mining      int32
	DBMutex     = &sync.Mutex{}
	UTXO_BUCKET = []byte("chainstate")
)
