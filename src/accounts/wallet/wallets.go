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

package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/YuriyLisovskiy/blockchain-go/src/config"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/secp256k1"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets(cfg config.Config) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile(cfg)
	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet
	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (ws Wallets) GetWallet(address string) (Wallet, error) {
	wallet, exists := ws.Wallets[address]
	if !exists {
		return Wallet{}, errors.New(fmt.Sprintf("wallet %s not found", address))
	}
	return *wallet, nil
}

func (ws *Wallets) LoadFromFile(cfg config.Config) error {
	if _, err := os.Stat(cfg.WalletsPath); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(cfg.WalletsPath)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(secp256k1.S256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}

func (ws Wallets) SaveToFile(cfg config.Config) {
	var content bytes.Buffer
	gob.Register(secp256k1.S256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(cfg.WalletsPath, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
