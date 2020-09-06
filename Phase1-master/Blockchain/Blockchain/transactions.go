package Blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct{
	ID []byte
	Inputs []TXInput
	Outputs []TXOutput
}

type TXOutput struct{
	Value int
	PubKey string
}

type TXInput struct{
	ID []byte
	Out int
	Sig string
}

func CoinbaseTx(to, data string) *Transaction{
	if data == ""{
		data = fmt.Sprintf("Coins to %s",to)
	}
	txin := TXInput{[]byte{}, -1,data}
	txout := TXOutput{100,to}

	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{txout}}
	tx.SetID()


	return &tx
}

func (tx *Transaction) SetID(){
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err:= encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool{
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TXInput) CanUnlock(data string) bool{
	return in.Sig == data
}

func (out *TXOutput)Canbeunlocked(data string)bool{
	return out.PubKey == data
}

func Newtransaction(from, to string, amount int, chain *BlockChain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error: Not sufficient funds")
	}

	for txid, outs :=  range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {

			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs,TXOutput{amount,to})

	if acc > amount{
		outputs = append(outputs,TXOutput{acc-amount,from})
	}
	tx := Transaction{nil,inputs,outputs}
	tx.SetID()
	return &tx
}
