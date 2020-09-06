# Phase2
For Phase2 of the csc 462 project

The blockchain comes with a commandline interface that is implemented in the main.go

In order to run it, simply open your command prompt and type

go run main.go createblockchain -address <Name for the address>

this would be the first step to do so as you would need to create a blockchain for your account

following that to check you account balance

go run main.go getbalance -address <Name of your address>

this would give the accounts available balance

if you want to send a transfer to another account

go run main.go send -from <from account name> -to <to account name> -amount <value to be sent>

finally in order to print the block chain,

go run main.go print

All the above commands generate several statistics from the badgerDB and are only there to further show the
performance of the implementation

In order to run this on your machine, several dependencies for the BadgerDB are necessary
	
	github.com/AndreasBriese/bbloom v0.0.0-20180913140656-343706a395b7 // indirect
	
	github.com/dgraph-io/badger v1.5.4
	
	github.com/dgryski/go-farm v0.0.0-20180109070241-2de33835d102 // indirect
	
	github.com/golang/protobuf v1.2.0 // indirect
	
	github.com/pkg/errors v0.8.0 // indirect
	
	golang.org/x/net v0.0.0-20181023162649-9b4f9f5ad519 // indirect
  
the code for the database will not function without these imports

instead of downloading them manually, the following link contains all the files necessary
for running the implementation and is made public

https://drive.google.com/drive/folders/1dPGpY1EELrhXjr-LJ7hIgKkEM4yuyioa?usp=sharing


  
Finally you need a file called tmp in the working directory to store the database
  

To run the web app, you must have a folder named "tmp" in the WebApp folder to store the blockchain.
After that, you can type "go test -run TestStartFresh" to start the web app, then go to localhost:8080 on your browser to visit it
Since golang tests have a 10 minute timeout, you can only test through this method for 10 minutes at a time before the test closes

Other network tests found in app_test.go
