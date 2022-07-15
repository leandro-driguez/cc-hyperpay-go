package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	Balance float32 `json:"Balance float32"`
	Bank    string  `json:"Bank string"`
}

// InitLedger adds a base set of accounts to the ledger
func (s *SmartContract) InitLedger(ctx CustomTransactionContextInterface) error {

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

// CreateAccount issues a new account to the world state with given details.
func (s *SmartContract) CreateAccount(ctx CustomTransactionContextInterface, id string, balance float32, bank string) error {

	exists := s.AccountExists(ctx, id)

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

func (s *SmartContract) ReadAccount(ctx CustomTransactionContextInterface, id string) (*Account, error) {

	contract := new(ReadonlySmartContract)
	account, err := contract.ReadAccount(ctx, id)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// DeleteAccount deletes an given account from the world state.
func (s *SmartContract) DeleteAccount(ctx CustomTransactionContextInterface, id string) error {

	exists := s.AccountExists(ctx, id)

	if !exists {
		return fmt.Errorf("the account %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AccountExists returns true when account with given ID exists in world state
func (s *SmartContract) AccountExists(ctx CustomTransactionContextInterface, id string) bool {

	accountJSON := ctx.GetData()

	return accountJSON != nil
}

func (s *SmartContract) Transfer(ctx CustomTransactionContextInterface, fromId, toId string, amount float32) error {
	// @todo q solo pueda hacer esto el duenyo d la cuenta fuente

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	contract := new(ReadonlySmartContract)

	fromAcc, err := contract.ReadAccount(ctx, fromId)
	if err != nil {
		return errors.New("source account doesn't exist")
	}

	toAcc, err := contract.ReadAccount(ctx, toId)
	if err != nil {
		return errors.New("source account doesn't exist")
	}

	fromAcc.Balance -= amount
	toAcc.Balance += amount

	fromAccJson, err := json.Marshal(fromAcc)
	if err != nil {
		return err
	}

	if err := ctx.GetStub().PutState(fromId, fromAccJson); err != nil {
		return err
	}

	toAccJson, err := json.Marshal(toAcc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(toId, toAccJson)
}

func beforeTransaction(ctx CustomTransactionContextInterface) error {

	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return errors.New("cannot get client's MSP-ID")
	}

	functionName, params := ctx.GetStub().GetFunctionAndParameters()

	if strings.Contains(functionName, "ReadAccount") {
		// aqui podemos hacer algo solo si se invoco la Tx ReadAccount
		return nil
	}

	contract := new(ReadonlySmartContract)
	account, err := contract.ReadAccount(ctx, params[0])

	if err != nil {
		return err
	}

	if (account.Bank == "JPMorgan Chase & Co." && mspID == "Org1MSP") ||
		(account.Bank == "Bank of America Corp." && mspID == "Org2MSP") {
		return errors.New("account does not belong to the executing org")
	}

	return nil
}

func afterTransaction(ctx CustomTransactionContextInterface, iface interface{}) error {
	return nil
}
