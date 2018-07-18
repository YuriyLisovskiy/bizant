// Copyright (c) 2018 Yuriy Lisovskiy
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

package blockchain

import (
	"fmt"
	"log"
	"time"
	"bytes"
	"math/big"
	"crypto/rand"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/hex"
	"crypto/sha256"
	"crypto/elliptic"
	w "github.com/YuriyLisovskiy/blockchain-go/src/wallet"
	"github.com/YuriyLisovskiy/blockchain-go/src/primitives"
)

type Transaction struct {
	ID        []byte
	VIn       []primitives.TXInput
	VOut      []primitives.TXOutput
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
	var inputs []primitives.TXInput
	var outputs []primitives.TXOutput
	for _, vin := range tx.VIn {
		inputs = append(inputs, primitives.TXInput{TxId: vin.TxId, VOut: vin.VOut, PubKey: nil, Signature: nil})
	}
	for _, vOut := range tx.VOut {
		outputs = append(outputs, primitives.TXOutput{Value: vOut.Value, PubKeyHash: vOut.PubKeyHash})
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

func NewCoinBaseTX(to string, fees float64, data string) Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}
		data = fmt.Sprintf("%x", randData)
	}
	txIn := primitives.TXInput{TxId: []byte{}, VOut: -1, Signature: nil, PubKey: []byte(data)}
	txOut := primitives.NewTXOutput(MINING_REWARD+fees, to)
	tx := Transaction{nil, []primitives.TXInput{txIn}, []primitives.TXOutput{*txOut}, time.Now().Unix(), 0}
	tx.ID = tx.Hash()
	return tx
}

func NewUTXOTransaction(wallet *w.Wallet, to string, amount, fee float64, utxoSet *UTXOSet) Transaction {
	var inputs []primitives.TXInput
	var outputs []primitives.TXOutput
	pubKeyHash := w.HashPubKey(wallet.PublicKey)
	acc, validOutputs := utxoSet.FindSpendableOutputs(pubKeyHash, amount,)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			inputs = append(inputs, primitives.TXInput{TxId: txID, VOut: out, Signature: nil, PubKey: wallet.PublicKey})
		}
	}
	from := fmt.Sprintf("%x", wallet.GetAddress())
	outputs = append(outputs, *primitives.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *primitives.NewTXOutput(acc-amount, from)) // a change
	}
	tx := Transaction{nil, inputs, outputs, time.Now().Unix(), 0}
	tx.ID = tx.Hash()
	tx.Fee = tx.CalculateFee(fee)
	return utxoSet.BlockChain.SignTransaction(tx, wallet.PrivateKey)
}

func (tx *Transaction) CalculateFee(feePerByte float64) float64 {
	if tx.IsCoinBase() {
		return 0.0
	}
	if feePerByte < MIN_FEE_PER_BYTE {
		feePerByte = MIN_FEE_PER_BYTE
	}
	return float64(len(tx.VIn)*148+len(tx.VOut)*34+10) * MIN_FEE_PER_BYTE
}

func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}
