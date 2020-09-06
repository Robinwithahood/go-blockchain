package main

import (
	"../go_blockchain/Blockchain"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

/*Adding a key-value pair using database BadgerDB, byte key-value pairs being stored
in folders

 */

type CommandLine struct{}

//Usage commands
func (cl *CommandLine) printUsage(){
	fmt.Println("Usage: ")
	fmt.Println(" getbalance -address <address> -get balance for an address")
	fmt.Println(" createblockchain -address <address> -creates a blockchain and sends an opening account reward of 100 to the address ")
	fmt.Println(" print -to print the chain")
	fmt.Println(" send -from <from_address> -to <to_address> -amount <Amount> -send the amount between addresses")
}

func (cl *CommandLine) validateArgs(){
	if len(os.Args)<2{
		cl.printUsage()
		runtime.Goexit() //exits the application by shutting down the routines, useful for badger db; prevents it from corrupting
	}
}

func (cl *CommandLine) Printchain() {
	chain := Blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.ChainIter()

	for {
		block := iter.Next()
		fmt.Printf("previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := Blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cl *CommandLine) createBlockChain(address string){
	chain := Blockchain.InitBlockchain(address)
	chain.Database.Close()
	fmt.Println("Completed creating a blockchain")
}

func (cl *CommandLine) getBalance(address string){
	chain := Blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs{
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n",address,balance)
}

func (cl *CommandLine) Send(from,to string, amount int){
	chain := Blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := Blockchain.Newtransaction(from,to,amount,chain)
	chain.AddBlock([]*Blockchain.Transaction{tx})
	fmt.Println("Finished sending amount")


}


func (cl *CommandLine) run(){
	cl.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	createBlockCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printcommand := flag.NewFlagSet("print",flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address","","The address to get balance is")
	createBlockchainAddress := createBlockCmd.String("address","","The address to create the blockchain for is")
	sendFrom := sendCmd.String("from","","From address")
	sendto := sendCmd.String("to","","To address")
	sendAmount := sendCmd.Int("amount",0,"Amount to send")

	switch os.Args[1]{
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "print":
		err := printcommand.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	default:
		cl.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed(){
		if *getBalanceAddress == ""{
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cl.getBalance(*getBalanceAddress)
	}

	if createBlockCmd.Parsed(){
		if *createBlockchainAddress == ""{
			createBlockCmd.Usage()
			runtime.Goexit()
		}
		cl.createBlockChain(*createBlockchainAddress)
	}

	if printcommand.Parsed(){
		cl.Printchain()
	}


	if sendCmd.Parsed(){
		if *sendFrom == "" || *sendto == "" || *sendAmount <= 0{
			sendCmd.Usage()
			runtime.Goexit()
		}
		cl.Send(*sendFrom,*sendto,*sendAmount)
	}
}

func main() {
	defer os.Exit(0) //only executes if go channel exits properly
	cl := CommandLine{}
	cl.run()
}
