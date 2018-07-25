// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import "fmt"

// TODO: Remove assertions before release.

// __assert__ will panic with a given formatted message if the given condition is false.
func __assert__(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: " + msg, v...))
	}
}
