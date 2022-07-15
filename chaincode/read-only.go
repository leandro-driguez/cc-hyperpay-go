package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Defines a private smart contract of read-only access
type ReadonlySmartContract struct {
	contractapi.Contract
}

// ReadAccount returns the account stored in the world state with given id.
func (s *ReadonlySmartContract) ReadAccount(ctx CustomTransactionContextInterface, id string) (*Account, error) {
	accountJSON := ctx.GetData()

	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	var err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
