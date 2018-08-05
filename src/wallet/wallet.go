// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the GNU General Public License v3.0 software license,
// see the accompanying file LICENSE or https://opensource.org/licenses/GPL-3.0

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
	"crypto/elliptic"
)

type Wallet struct {
	PrivateKey []byte
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

func newKeyPair() (publicKey []byte, privateKey []byte) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	publicKey = elliptic.Marshal(secp256k1.S256(), key.X, key.Y)
	privateKey = make([]byte, 32)
	blob := key.D.Bytes()
	copy(privateKey[32-len(blob):], blob)
	return
}
