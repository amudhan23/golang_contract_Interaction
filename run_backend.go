package main

import(
	"fmt"
	"log"
  "crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core"
  "github.com/ethereum/go-ethereum/common"
	"math/big"
   import_contracts "./contracts"
)


func main(){

		authSender,_ := AuthAndAddressGeneration()
		alloc := make(core.GenesisAlloc)
		alloc[authSender.From] = core.GenesisAccount{Balance: big.NewInt(11000000000000000)}  //1.1 ether

		//client
		client := backends.NewSimulatedBackend(alloc,600000000000000000)

		contractAddress,_,_, err:=import_contracts.DeploySimpleStorage(
			authSender,
			client,
		)
		if err!=nil{
			log.Fatal(err)
		}
		client.Commit()

		instance,err:= import_contracts.NewSimpleStorage(contractAddress,client)
		if err != nil {
				log.Fatal(err)
		}

		StoredDataValue,err:=instance.StoredData(nil)
		if err != nil {
				log.Fatal(err)
		}

		fmt.Println("Stored Data Value before setting : ",StoredDataValue)


		_,err=instance.Set(&bind.TransactOpts{
				From:authSender.From,
				Signer:authSender.Signer,
				Value: nil,
		},big.NewInt(8))
		if err != nil {
				log.Fatal(err)
		}
		client.Commit()

		StoredDataValue,err=instance.StoredData(nil)
		if err != nil {
				log.Fatal(err)
		}

		fmt.Println("Stored Data Value after setting : ",StoredDataValue)

}


//generating a  computing party in the blockchain
func AuthAndAddressGeneration() (*bind.TransactOpts,common.Address){

	privateKey,err:=crypto.GenerateKey()
  if err!=nil{
    log.Fatal(err)
  }
  auth:=bind.NewKeyedTransactor(privateKey)

  publicKey := privateKey.Public()
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
      log.Fatal("error casting public key to ECDSA")
  }
  address := crypto.PubkeyToAddress(*publicKeyECDSA)
  return auth,address
}
