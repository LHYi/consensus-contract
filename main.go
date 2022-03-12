package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	responsecredit "github.com/LHYi/consensus-contract/consensus"
)

func main() {

	contract := new(responsecredit.Contract)
	contract.TransactionContextHandler = new(responsecredit.TransactionContext)
	contract.Name = "org.creditnet.responsecredit"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "ResponseCreditChaincode"
	chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}
