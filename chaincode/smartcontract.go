package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Account
type SmartContract struct {
	contractapi.Contract
}

// Account describes basic details of what makes up a simple account
//Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
type Account struct {
	ID      string  `json:"ID"`
	Balance float32 `json:"PatientName float32"`
	Bank    string  `json:"Bank string"`
}

// InitLedger adds a base set of accounts to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	accounts := []Account{
		{ID: "account1", Balance: 100, Bank: "JPMorgan Chase & Co."},
		{ID: "account2", Balance: 200, Bank: "Bank of America Corp."},
		{ID: "account3", Balance: 300, Bank: "JPMorgan Chase & Co."},
		{ID: "account4", Balance: 400, Bank: "Bank of America Corp."},
		{ID: "account5", Balance: 500, Bank: "JPMorgan Chase & Co."},
	}

	for _, account := range accounts {
		accountJSON, err := json.Marshal(account)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(account.ID, accountJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ReadAccount returns the account stored in the world state with given id.
func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, id string) (*Account, error) {
	accountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// CreateAccount issues a new account to the world state with given details.
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id string, balance float32, bank string) error {
	// Checking if the tx is being executed by org1
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return errors.New("cannot get client's MSP-ID")
	}
	if !(mspID == "Org1MSP" && bank == "JPMorgan Chase & Co." ||
		mspID == "Org2MSP" && bank == "Bank of America Corp.") {
		return fmt.Errorf("you have no access to this Tx")
	}

	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the account %s already exists", id)
	}

	account := Account{
		ID:      id,
		Balance: balance,
		Bank:    bank,
	}
	accountJSON, err := json.Marshal(account)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, accountJSON)
}

// DeleteAccount deletes an given account from the world state.
func (s *SmartContract) DeleteAccount(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AccountExists(ctx, id)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the account %s does not exist", id)
	}

	//Here we check if the account belongs to the bank that is trying to delete it
	account_existing, err := s.ReadAccount(ctx, id)
	if err != nil {
		return err
	}
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return errors.New("cannot get client's MSP-ID")
	}
	if (account_existing.Bank == "JPMorgan Chase & Co." && mspID == "Org1MSP") ||
		(account_existing.Bank == "Bank of America Corp." && mspID == "Org2MSP") {
		return errors.New("account does not belong to the executing org")
	}

	return ctx.GetStub().DelState(id)
}

// AccountExists returns true when account with given ID exists in world state
func (s *SmartContract) AccountExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	accountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return accountJSON != nil, nil
}
