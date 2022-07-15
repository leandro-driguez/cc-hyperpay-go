/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/lllrdgz/cc-hyperpay-go/hyperpay-transfer/chaincode"
)

func main() {
	smartContract := new(chaincode.SmartContract)
	smartContract.TransactionContextHandler = new(chaincode.CustomTransactionContext)
	smartContract.BeforeTransaction = chaincode.GetWorldState
	smartContract.UnknownTransaction = chaincode.UnknownTransactionHandler

	accountChaincode, err := contractapi.NewChaincode(smartContract)

	if err != nil {
		log.Panicf("Error creating account-transfer-basic chaincode: %v", err)
	}

	if err := accountChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
