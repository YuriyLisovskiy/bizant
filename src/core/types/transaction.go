// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package types

import (
	"fmt"
	"log"
	"bytes"
	"math/big"
	"crypto/rand"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/hex"
	"crypto/sha256"
	"crypto/elliptic"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/types/tx_io"
)

type Transaction struct {
	ID        []byte
	VIn       []tx_io.TXInput
	VOut      []tx_io.TXOutput
	Timestamp int64
	Fee       float64
}

func (tx Transaction) IsCoinBase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxId) == 0 && tx.VIn[0].VOut == -1
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) Transaction {
	if tx.IsCoinBase() {
		return *tx
	}
	for _, vin := range tx.VIn {
		if prevTXs[hex.EncodeToString(vin.TxId)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inID, vIn := range txCopy.VIn {
		prevTx := prevTXs[hex.EncodeToString(vIn.TxId)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTx.VOut[vIn.VOut].PubKeyHash
		//	dataToSign := fmt.Sprintf("%x\n", txCopy)
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.Serialize())
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.VIn[inID].Signature = signature
		txCopy.VIn[inID].PubKey = nil
	}
	return *tx
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []tx_io.TXInput
	var outputs []tx_io.TXOutput
	for _, vin := range tx.VIn {
		inputs = append(inputs, tx_io.TXInput{TxId: vin.TxId, VOut: vin.VOut, PubKey: nil, Signature: nil})
	}
	for _, vOut := range tx.VOut {
		outputs = append(outputs, tx_io.TXOutput{Value: vOut.Value, PubKeyHash: vOut.PubKeyHash})
	}
	txCopy := Transaction{ID: tx.ID, VIn: inputs, VOut: outputs, Timestamp: tx.Timestamp, Fee: tx.Fee}
	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	if len(tx.VIn) == 0 {
		log.Panic("ERROR: bad-txns-vin-empty")
	}
	if len(tx.VOut) == 0 {
		log.Panic("ERROR: bad-txns-vout-empty")
	}
	for _, vin := range tx.VIn {
		if prevTXs[hex.EncodeToString(vin.TxId)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()

	fmt.Printf("LOG: tx hash: %x == tx copy hash: %x ? %t",
		tx.Hash(), txCopy.Hash(),
		fmt.Sprintf("%x", tx.Hash()) == fmt.Sprintf("%x", txCopy.Hash()),
	)

	curve := elliptic.P256()
	for inID, vin := range tx.VIn {
		prevTx := prevTXs[hex.EncodeToString(vin.TxId)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTx.VOut[vin.VOut].PubKeyHash
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		//		dataToVerify := fmt.Sprintf("%x\n", txCopy)
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txCopy.Serialize(), &r, &s) == false {
			return false
		}
		txCopy.VIn[inID].PubKey = nil
	}
	return true
}

func (tx *Transaction) CalculateFee(feePerByte float64) float64 {
	if tx.IsCoinBase() {
		return 0.0
	}
	if feePerByte < vars.MIN_FEE_PER_BYTE {
		feePerByte = vars.MIN_FEE_PER_BYTE
	}
	return float64(len(tx.VIn)*148+len(tx.VOut)*34+10) * vars.MIN_FEE_PER_BYTE
}
