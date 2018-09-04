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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

func TestConfig_Default(t *testing.T) {
	ip, err := getIp()
	if err != nil {
		t.Errorf("config.TestConfig_Defualt: can't test config")
	}
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Errorf("config.TestConfig_Defualt: can't test config")
	}
	expected := Config{
		Ip: ip,
		Port: 8000,
		ChainPath: absPath + "/data/" + fmt.Sprintf(utils.DBFile, 8000),
		WalletsPath: absPath + "/data/" + fmt.Sprintf(utils.WalletFile, 8000),
	}
	actual, err := Default()
	if err != nil {
		t.Errorf("config.TestConfig_Defualt: can't test config")
	}
	if actual.Ip != expected.Ip {
		t.Errorf("config.TestConfig_Defualt, Ip: %s != %s", actual.Ip, expected.Ip)
	}
	if actual.Port != expected.Port {
		t.Errorf("config.TestConfig_Defualt, Port: %d != %d", actual.Port, expected.Port)
	}
	if actual.WalletsPath != expected.WalletsPath {
		t.Errorf("config.TestConfig_Defualt, WalletsPath: %s != %s", actual.WalletsPath, expected.WalletsPath)
	}
	if actual.ChainPath != expected.ChainPath {
		t.Errorf("config.TestConfig_Defualt, ChainPath: %s != %s", actual.ChainPath, expected.ChainPath)
	}
}

func TestConfig_SetIp(t *testing.T) {
	cfg := Config{}
	cfg = cfg.SetIp("127.0.0.1")
	if cfg.Ip != "127.0.0.1" {
		t.Errorf("config.TestConfig_SetIp: %s != %s", cfg.Ip, "127.0.0.1")
	}
}

func TestConfig_SetPort(t *testing.T) {
	cfg := Config{}
	cfg = cfg.SetPort(3000)
	if cfg.Port != 3000 {
		t.Errorf("config.TestConfig_SetPort: %d != %d", cfg.Port, 3000)
	}
}

func TestConfig_SetChainPath(t *testing.T) {
	cfg := Config{}
	cfg = cfg.SetChainPath("some/path/to/chain")
	if cfg.ChainPath != "some/path/to/chain" {
		t.Errorf("config.TestConfig_SetChainPath: %s != %s", cfg.ChainPath, "some/path/to/chain")
	}
}

func TestConfig_SetWalletsPath(t *testing.T) {
	cfg := Config{}
	cfg = cfg.SetWalletsPath("some/wallets/path")
	if cfg.WalletsPath != "some/wallets/path" {
		t.Errorf("config.TestConfig_SetWalletsPath: %s != %s", cfg.WalletsPath, "some/wallets/path")
	}
}

func TestConfig_Exists(t *testing.T) {
	cfg := Config{}
	exists := Exists()
	if exists != false {
		t.Errorf("config.TestConfig_Exists: %t != %t", exists, false)
	}
	err := cfg.Save()
	if err != nil {
		panic(err)
	}
	exists = Exists()
	if exists != true {
		t.Errorf("config.TestConfig_Exists: %t != %t", exists, true)
	}
	os.Remove(configLocation)
}

func TestConfig_Load(t *testing.T) {
	expected, err := Default()
	if err != nil {
		t.Errorf("config.TestConfig_Load: can't test config")
	}
	err = expected.Save()
	if err != nil {
		t.Errorf("config.TestConfig_Load: can't test config")
	}
	actual, err := LoadConfig()
	if err != nil {
		t.Errorf("config.TestConfig_Load: can't test config")
	}
	os.Remove(configLocation)
	if actual.Ip != expected.Ip {
		t.Errorf("config.TestConfig_Load, Ip: %s != %s", actual.Ip, expected.Ip)
	}
	if actual.Port != expected.Port {
		t.Errorf("config.TestConfig_Load, Port: %d != %d", actual.Port, expected.Port)
	}
	if actual.WalletsPath != expected.WalletsPath {
		t.Errorf("config.TestConfig_Load, WalletsPath: %s != %s", actual.WalletsPath, expected.WalletsPath)
	}
	if actual.ChainPath != expected.ChainPath {
		t.Errorf("config.TestConfig_Load, ChainPath: %s != %s", actual.ChainPath, expected.ChainPath)
	}
}
