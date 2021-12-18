package main

import (
	"fmt"

	"IFC/inventoryfinancingpaper"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	contract := new(inventoryfinancingpaper.Contract)
	contract.TransactionContextHandler = new(inventoryfinancingpaper.TransactionContext)
	contract.Name = "org.papernet.inventoryfinancingpaper"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "InventoryFinancingPaperChaincode"
	chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}
