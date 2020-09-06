package Blockchain


/* The file is meant to act as a simpler consensus algorithm, implementing the Proof of Work algorithm
the algorithm, takes the data from the block and creates a counter. It then validates the proof against a set
of requirements specified. The notion that changing a block inside a chain requires recalculating hashes for the block
itself and the every block before the target block, makes it a secure yet simple algorithm.

 */
import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 18 //A measure of requirements, which is generally intended to increase with time and users but static in our case

type PoW struct{
	Block *Block
	Target *big.Int
}

//The method takes a block and pairs with the target


func NewProof (b*Block) *PoW{
	target := big.NewInt(1)
	target.Lsh(target,uint(256-Difficulty))//number of bytes in a hash - number of 0's and shift left

	proof := &PoW{b,target}
	return proof
}

func (proof *PoW) InitData(nonce int)[]byte{
	//combine the data in order to push it onto the block
	data := bytes.Join([][] byte{
		proof.Block.PrevHash,
		proof.Block.HashTransactions(),
		Hexa(int64(nonce)),
		Hexa(int64(Difficulty)),
	},
	[]byte{})
	return data
}

//Utility method that converts the data into a big endian hexadecimal format
func Hexa(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}

func (proof *PoW) Run() (int, []byte){
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	//essentially an infinite loop, preparing the data, creating a hash and comparing it with target hash
	for nonce < math.MaxInt64 {
		data := proof.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(proof.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}
/* The validation method is meant to validate the NewProof() method, essentially
validates it against the specified rules. In this case there must be a preceding set
of few 0s in the first few bytes.
 */
func (proof *PoW) Validate() bool{
	var intHash big.Int

	data := proof.InitData(proof.Block.Nonce)

	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(proof.Target) == -1
}
