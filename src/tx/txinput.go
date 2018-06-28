package tx

import (
	"bytes"
	"github.com/YuriyLisovskiy/blockchain-go/src/wallet"
)

type TXInput struct {
	TxId      []byte
	VOut      int
	Signature []byte
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	return bytes.Compare(wallet.HashPubKey(in.PubKey), pubKeyHash) == 0
}
