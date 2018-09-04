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

package jwt

import (
	"testing"
	"time"
)

func TestNewClaims(t *testing.T) {
	claims := NewClaims()
	if !claims.Contains("iat") {
		t.Errorf("jwt.TestNewClaims: claims is not valid")
	}
}

var Set_TestClaimsData = []struct {
	key      string
	val      string
	expected string
}{
	{
		key: "key1",
		val: "val1",
		expected: "val1",
	},
	{
		key: "key2",
		val: "val2",
		expected: "val2",
	},
	{
		key: "key3",
		val: "val3",
		expected: "val3",
	},
	{
		key: "key4",
		val: "val4",
		expected: "val4",
	},{
		key: "key5",
		val: "val5",
		expected: "val5",
	},
	{
		key: "key6",
		val: "val6",
		expected: "val6",
	},{
		key: "key7",
		val: "val7",
		expected: "val7",
	},
	{
		key: "key8",
		val: "val8",
		expected: "val8",
	},
}

func TestClaims_Set(t *testing.T) {
	claims := NewClaims()
	for _, data := range Set_TestClaimsData {
		claims.Set(data.key, data.val)
	}
	for _, data := range Set_TestClaimsData {
		str, _ := claims.GetString(data.key)
		if str != data.expected {
			t.Errorf("jwt.TestClaims_Set: %s != %s", str, data.expected)
		}
	}
}

func TestClaims_SetTime(t *testing.T) {
	claims := NewClaims()
	tm := time.Now()
	claims.SetTime("exp", tm)
	exp, _ := claims.GetTime("exp")
	if exp.Unix() != tm.Unix() {
		t.Errorf("jwt.TestClaims_SetTime: %d != %d", exp.Unix(), tm.Unix())
	}
}
