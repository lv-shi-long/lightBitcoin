package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type TXInput struct {
	TXID   []byte
	Index  int64
	Addres string
}

type TXOutput struct {
	Value  float64
	Addres string
}

type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

// 设置交易的ID，先把 本身的ID 置为Nil,再进行计算。
func (tx *Transaction) SetTransactionID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	encoder.Encode(tx)
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

// 挖矿交易，没有输入，只有输出，默认设置 12.5个比特币
func NewCoinBase(miner string) *Transaction {

	input := TXInput{
		TXID:   nil,
		Index:  -1,
		Addres: miner + "dig out",
	}
	output := TXOutput{
		Value:  12.5,
		Addres: miner,
	}
	tx := Transaction{
		TXID:      nil,
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	tx.SetTransactionID()
	return &tx
}
