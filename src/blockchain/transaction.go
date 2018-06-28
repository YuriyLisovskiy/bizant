package blockchain

import (
	"fmt"
	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	txPkg "github.com/YuriyLisovskiy/blockchain-go/src/tx"
)

type Transaction struct {
	ID   []byte
	VIn  []txPkg.TXInput
	VOut []txPkg.TXOutput
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txIn := txPkg.TXInput{[]byte{}, -1, data}
	txOut := txPkg.TXOutput{utils.Reward, to}
	tx := Transaction{nil, []txPkg.TXInput{txIn}, []txPkg.TXOutput{txOut}}
	tx.SetID()
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []txPkg.TXInput
	var outputs []txPkg.TXOutput
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
		if err != nil {
			panic(err)
		}
		for _, out := range outs {
			input := txPkg.TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, txPkg.TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, txPkg.TXOutput{acc - amount, from}) // a change
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

func (tx Transaction) IsCoinBase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxId) == 0 && tx.VIn[0].VOut == -1
}