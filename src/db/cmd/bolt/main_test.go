// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package db

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	. "github.com/boltdb/bolt/cmd/bolt"
)

// open creates and opens a Bolt database in the temp directory.
func open(fn func(*bolt.DB)) {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())
	defer os.RemoveAll(f.Name())

	db, err := bolt.Open(f.Name(), 0600)
	if err != nil {
		panic("db open error: " + err.Error())
	}
	fn(db)
}

// run executes a command against the CLI and returns the output.
func run(args ...string) string {
	args = append([]string{"bolt"}, args...)
	NewApp().Run(args)
	return strings.TrimSpace(LogBuffer())
}
