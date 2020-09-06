package Blockchain

import (
	"../github.com/dgraph-io/badger"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
)

const dbpath = "./tmp/blocks"
const dbFile = "./tmp/blocks/MANIFEST"
const genesis = "The first block"

type BlockChain struct{
	LastHash []byte
	Database  *badger.DB
}

//structure to iterate through the blockchain
type Iterator struct{
	currenth []byte
	DB *badger.DB
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lasthash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh")) //"lh" stands for lasthash, gives a pointer to item stored in the lh key and an error
		Handle(err)
		var hash []byte
		err = item.Value(func(val []byte) error {
			hash = append([]byte{}, val...)
			return nil
		})
		lasthash = hash
		return err
	})
	Handle(err)

	newblock := CreateBlock(transactions,lasthash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newblock.Hash,newblock.serial())
		Handle(err)
		err = txn.Set([]byte("lh"),newblock.Hash)

		chain.LastHash = newblock.Hash
		return err
	})
	Handle(err)
}

/*Initializing Blockchain using BadgerDB.

 */

func DBexists() bool{
	if _,err := os.Stat(dbFile); os.IsNotExist(err){
		return false
	}
	return true
}

func InitBlockchain(address string) *BlockChain{
	var lastHash []byte

	if DBexists(){
		fmt.Println("BlockChain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbpath)
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)
	Handle(err)


	//two ways to access the database, view or update, update allows for writes
	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesis)
		gene := Genesis(cbtx)
		fmt.Println("Genesis Created")
		err = txn.Set(gene.Hash, gene.serial())
		Handle(err)
		err = txn.Set([]byte("lh"), gene.Hash)
		lastHash = gene.Hash
		return err
	})

	Handle(err)
	chain := BlockChain{lastHash,db}
	return &chain
}

func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No Existing Blockchain found")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbpath)
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		var Hash []byte
		err = item.Value(func(val []byte) error {
			Hash = append([]byte{}, val...)
			return nil
		})
		lastHash = Hash
		return err
	})
	Handle(err)
	chain := BlockChain{lastHash, db}
	return &chain
}


//Manually implementing an iterator instead of pre-built functions, iterating backwards as we start from the last hash
func (chain *BlockChain) ChainIter() *Iterator{
	iter := &Iterator{chain.LastHash,chain.Database}
	return iter
}

//Function that returns a pointer to the next block
func (iter *Iterator) Next() *Block{
	var block *Block

	//in order to find the next block, we read through the chain and hence use the view function instead of update
	err := iter.DB.View(func(txn *badger.Txn) error {
		item,err := txn.Get(iter.currenth)
		Handle(err)

		var encblock []byte
		err = item.Value(func(val []byte) error {
			encblock = append([]byte{},val...)
			return nil
		})
		block = deserial(encblock)

		return err
	})
	Handle(err)

	iter.currenth = block.PrevHash

	return block
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.ChainIter()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIDx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIDx{
							continue Outputs
						}
					}
				}
				if out.Canbeunlocked(address){
					unspentTxs = append(unspentTxs,*tx)
				}
			}
			if tx.IsCoinbase() == false{
				for _,in := range tx.Inputs{
					if in.CanUnlock(address){
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Out)
					}
				}
			}
		}
		if len(block.PrevHash) == 0{
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindUTXO(address string) []TXOutput{
	var UTXOs []TXOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions{
		for _,out := range tx.Outputs{
			if out.Canbeunlocked(address){
				UTXOs = append(UTXOs,out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int){
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

	Work:
		for _, tx := range unspentTxs{
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Outputs{
				if out.Canbeunlocked(address) && accumulated < amount{
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID],outIdx)

					if accumulated >= amount{
						break Work
					}
				}
			}
		}
		return accumulated,unspentOuts
}