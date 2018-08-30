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

package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"

	"github.com/YuriyLisovskiy/blockchain-go/src/core/types/tx_io"
	"github.com/YuriyLisovskiy/blockchain-go/src/core/vars"
	"github.com/YuriyLisovskiy/blockchain-go/src/crypto/secp256k1"
)

type Transaction struct {
	Hash      []byte
	VIn       []tx_io.TXInput
	VOut      []tx_io.TXOutput
	Timestamp int64
	Fee       float64
}

func (tx Transaction) IsCoinBase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].PreviousTx) == 0 && tx.VIn[0].VOut == -1
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

func (tx *Transaction) CalcHash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.Hash = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Sign(privateKey []byte, prevTXs map[string]Transaction) Transaction {
	if tx.IsCoinBase() {
		return *tx
	}
	for _, vin := range tx.VIn {
		if prevTXs[hex.EncodeToString(vin.PreviousTx)].Hash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	//	fmt.Printf("\n\nPUB KEY (sign): %x\n", append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...))

	txCopy := tx.TrimmedCopy()

	for inID, vIn := range tx.VIn {
		prevTx := prevTXs[hex.EncodeToString(vIn.PreviousTx)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTx.VOut[vIn.VOut].PubKeyHash
		//	dataToSign := fmt.Sprintf("%x\n", txCopy)
		signature, err := secp256k1.Sign(tx.Hash, privateKey)
		if err != nil {
			log.Panic(err)
		}
		tx.VIn[inID].Signature = signature
		//	txCopy.VIn[inID].PubKey = nil
	}
	return *tx
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []tx_io.TXInput
	var outputs []tx_io.TXOutput
	for _, vin := range tx.VIn {
		inputs = append(inputs, tx_io.TXInput{PreviousTx: vin.PreviousTx, VOut: vin.VOut, PubKey: nil, Signature: nil})
	}
	for _, vOut := range tx.VOut {
		outputs = append(outputs, tx_io.TXOutput{Value: vOut.Value, PubKeyHash: vOut.PubKeyHash})
	}
	txCopy := Transaction{Hash: tx.Hash, VIn: inputs, VOut: outputs, Timestamp: tx.Timestamp, Fee: tx.Fee}
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
		if prevTXs[hex.EncodeToString(vin.PreviousTx)].Hash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()

	//	fmt.Printf("LOG: tx hash: %x == tx copy hash: %x ? %t",
	//		tx.CalcHash(), txCopy.CalcHash(),
	//		fmt.Sprintf("%x", tx.CalcHash()) == fmt.Sprintf("%x", txCopy.CalcHash()),
	//	)
	for inID, vin := range tx.VIn {
		prevTx := prevTXs[hex.EncodeToString(vin.PreviousTx)]
		//	txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTx.VOut[vin.VOut].PubKeyHash
		/*
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		*/
		//	x, y := big.Int{}, big.Int{}
		//	x.SetBytes(vin.PubKey[:(len(vin.PubKey) / 2)])
		//	y.SetBytes(vin.PubKey[(len(vin.PubKey) / 2):])

		//	pubKey := ecdsa.PublicKey{X: &x, Y: &y, Curve: secp256k1.S256()}

		//	fmt.Printf("\nVIN PUB KEY (verify): %x\n\n", vin.PubKey)
		//	fmt.Printf("\nPUB KEY (verify): %x\n\n\n", pubKey)

		if !secp256k1.VerifySignature(vin.PubKey, tx.Hash, vin.Signature) {
			return false
		}
		//	txCopy.VIn[inID].PubKey = nil
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
