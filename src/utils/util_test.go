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

package utils

import (
	"os"
	"testing"
)

func TestDBExists(test *testing.T) {
	dbFile := "some_db.db"

	actual := DBExists(dbFile)
	if actual != false {
		test.Errorf("utils.DBExists: %t != %t", actual, false)
	}

	file, err := os.Create(dbFile)
	if err != nil {
		test.Error(err)
	}
	defer file.Close()

	actual = DBExists(dbFile)
	if actual != true {
		test.Errorf("utils.DBExists: %t != %t", actual, true)
	}

	err = os.Remove(dbFile)
	if err != nil {
		test.Error(err)
	}

}

var IntToHex_Data = []struct{
	input int64
	expected []byte
} {
	{
		input: 10,
		expected: []byte{0, 0, 0, 0, 0, 0, 0, 10},
	},
	{
		input: 0,
		expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		input: 234567897654,
		expected: []byte{0, 0, 0, 54, 157, 86, 18, 54},
	},
	{
		input: 74562,
		expected: []byte{0, 0, 0, 0, 0, 1, 35, 66},
	},
	{
		input: 190765,
		expected: []byte{0, 0, 0, 0, 0, 2, 233, 45},
	},
}

func TestIntToHex(test *testing.T) {
	for j, data := range IntToHex_Data {
		actual := IntToHex(data.input)
		if len(actual) != len(data.expected) {
			test.Errorf("utils.TestIntToHex[%d], size: %d != %d", j, len(actual), len(data.expected))
		}
		for i, b := range actual {
			if b != data.expected[i] {
				test.Errorf("utils.TestIntToHex[%d], data: %s != %s", j, actual, data.expected)
			}
		}
	}
}

var ReverseBytes_Data = []struct{
	input []byte
	expected []byte
} {
	{
		input: []byte{0, 0, 0, 0, 0, 0, 0, 10},
		expected: []byte{10, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 54, 157, 86, 18, 54},
		expected: []byte{54, 18, 86, 157, 54, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 0, 0, 1, 35, 66},
		expected: []byte{66, 35, 1, 0, 0, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 0, 0, 2, 233, 45},
		expected: []byte{45, 233, 2, 0, 0, 0, 0, 0},
	},
}

func TestReverseBytes(test *testing.T) {
	for j, data := range ReverseBytes_Data {
		actual := data.input
		ReverseBytes(actual)
		if len(actual) != len(data.expected) {
			test.Errorf("utils.TestReverseBytes[%d], size: %d != %d", j, len(actual), len(data.expected))
		}
		for i, b := range actual {
			if b != data.expected[i] {
				test.Errorf("utils.TestReverseBytes[%d], data: %s != %s", j, actual, data.expected)
			}
		}
	}
}
