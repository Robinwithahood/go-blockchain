package Blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct{
	Hash []byte
	Transactions []*Transaction
	PrevHash []byte
	Nonce int
}



func CreateBlock(txs []*Transaction, prevHash []byte)*Block{
	block := &Block{[]byte{},txs,prevHash,0}
	proof := NewProof(block)
	nonce,hash := proof.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}



func Genesis(coinbase *Transaction) *Block{
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}


/*Need to serialize and deserialize, blocks for the database

 */

func (block *Block) serial() []byte{
	var res bytes.Buffer
	enc := gob.NewEncoder(&res)
	err := enc.Encode(block)

	Handle(err)

	return res.Bytes()

}

//outputs a pointer to the block

func deserial(data []byte) *Block{
	var block Block

	dec := gob.NewDecoder(bytes.NewReader(data))

	err := dec.Decode(&block)

	Handle(err)
	return &block
}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}

func (b *Block) HashTransactions() []byte{

	var txHashes [][]byte
	var txHash [32]byte

	for _,tx := range b.Transactions{
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))

	return txHash[:]
}