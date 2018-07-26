// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package wallet

import (
	"log"
	"bytes"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"

	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"github.com/YuriyLisovskiy/blockchain-go/src/secp256k1"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{WALLLET_VERSION}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)
	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-ADDRESS_CHECKSUM_LEN:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1: len(pubKeyHash)-ADDRESS_CHECKSUM_LEN]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:ADDRESS_CHECKSUM_LEN]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := secp256k1.S256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
