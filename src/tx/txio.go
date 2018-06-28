package tx

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

type TXInput struct {
	TxId      []byte
	VOut      int
	ScriptSig string
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
