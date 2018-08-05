// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

package vars

import "sync"

var (
	Syncing      int32
	DBMutex     = &sync.Mutex{}
	UTXO_BUCKET = []byte("chainstate")
)
