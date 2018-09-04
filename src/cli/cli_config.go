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

package cli

import "github.com/YuriyLisovskiy/blockchain-go/src/config"

func (cli *CLI) setConfig(ip string, port int, chainPath, walletsPath string) error {
	cfg := config.Config{}
	var err error
	if config.Exists() {
		cfg, err = config.LoadConfig()
		if err != nil {
			return err
		}
	}
	if ip != "" {
		cfg = cfg.SetIp(ip)
	}
	if port != -1 {
		cfg = cfg.SetPort(port)
	}
	if chainPath != "" {
		cfg = cfg.SetChainPath(chainPath)
	}
	if walletsPath != "" {
		cfg = cfg.SetWalletsPath(walletsPath)
	}
	return cfg.Save()
}

func (cli *CLI) setDefaultConfig() error {
	cfg, err := config.Default()
	if err != nil {
		return err
	}
	return cfg.Save()
}
