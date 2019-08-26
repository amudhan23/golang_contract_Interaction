package main

import(
	"fmt"
	"log"
	"io"
	"os"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
   import_contracts "./contracts"
)


func main(){

		var nodeReader io.Reader
		nodeReader, err := os.Open("/home/amudhan/workspace1/node1/keystore/UTC--2019-07-25T15-02-05.471118987Z--bfbd4c070743dfa221f2f5d704dd8bd0f045e730")
		if err != nil{
				log.Fatal(err)
		}

		authSender,err	:= bind.NewTransactor(nodeReader,"qwerty")
		if err != nil{
				log.Fatal(err)
		}


		//client
		client,err := ethclient.Dial("/home/amudhan/workspace1/node1/geth.ipc")
		if err != nil {
					log.Fatal(err)
		}

		ctx := context.Background()

		_,tx,_, err:=import_contracts.DeploySimpleStorage(
			authSender,
			client,
		)
		if err!=nil{
			log.Fatal(err)
		}

		contractAddress,err:=bind.WaitDeployed(ctx,client,tx)
		if err!=nil{
			log.Fatal(err)
		}


		instance,err:= import_contracts.NewSimpleStorage(contractAddress,client)
		if err != nil {
				log.Fatal(err)
		}

		StoredDataValue,err:=instance.StoredData(nil)
		if err != nil {
				log.Fatal(err)
		}

		fmt.Println("Stored Data Value before setting : ",StoredDataValue)




		tx,err=instance.Set(&bind.TransactOpts{
				From:authSender.From,
				Signer:authSender.Signer,
				Value: nil,
		},big.NewInt(8))
		if err != nil {
				log.Fatal(err)
		}

		_,err=bind.WaitMined(ctx,client,tx)
		if err!=nil{
			log.Fatal(err)
		}



		StoredDataValue,err=instance.StoredData(nil)
		if err != nil {
				log.Fatal(err)
		}

		fmt.Println("Stored Data Value after setting : ",StoredDataValue)

}
