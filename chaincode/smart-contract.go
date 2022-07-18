package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Account
type SmartContract struct {
	contractapi.Contract
}

// Account describes basic details of what makes up a simple account
type Account struct {
	ID      string  `json:"ID"`
	Balance float32 `json:"Balance"`
	Bank    string  `json:"Bank"`
}

// TxRecord structure used to return the transaction history result of an account
type TxRecord struct {
	Record    *Account  `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
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

	// For each account encoding and save it
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

// AccountExists returns true when account with given ID exists in world state
func (s *SmartContract) AccountExists(ctx contractapi.TransactionContextInterface, accountID string) (bool, error) {

	// Get client org id and verify it matches peer org id.
	clientOrgID, err := getClientOrgID(ctx)
	if err != nil {
		return false, err
	}
	err = verifyClientOrgMatchesPeerOrg(clientOrgID)
	if err != nil {
		return false, err
	}

	accountJSON, err := ctx.GetStub().GetState(accountID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return accountJSON != nil, nil
}

// ReadAccount returns the account stored in the world state with given id.
func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, accountID string) (*Account, error) {

	// Get client org id and verify it matches peer org id.
	clientOrgID, err := getClientOrgID(ctx)
	if err != nil {
		return nil, err
	}
	err = verifyClientOrgMatchesPeerOrg(clientOrgID)
	if err != nil {
		return nil, err
	}

	accountJSON, err := ctx.GetStub().GetState(accountID)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", accountID)
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

	// Get client org id and verify it matches peer org id.
	clientOrgID, err := getClientOrgID(ctx)
	if err != nil {
		return err
	}
	err = verifyClientOrgMatchesPeerOrg(clientOrgID)
	if err != nil {
		return err
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

	err = ctx.GetStub().PutState(account.ID, accountJSON)
	if err != nil {
		return fmt.Errorf("failed to put asset in public data: %v", err)
	}

	// Set the endorsement policy such that an owner org peer is required to endorse future updates
	endorsingOrgs := []string{clientOrgID}
	err = setAssetStateBasedEndorsement(ctx, account.ID, endorsingOrgs)
	if err != nil {
		return fmt.Errorf("failed setting state based endorsement for buyer and seller: %v", err)
	}

	return nil
}

// DeleteAccount deletes an given account from the world state.
func (s *SmartContract) DeleteAccount(ctx contractapi.TransactionContextInterface, accountID string) error {

	// Get client org id and verify it matches peer org id.
	clientOrgID, err := getClientOrgID(ctx)
	if err != nil {
		return err
	}
	err = verifyClientOrgMatchesPeerOrg(clientOrgID)
	if err != nil {
		return err
	}

	exists, err := s.AccountExists(ctx, accountID)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the account %s does not exist", accountID)
	}

	return ctx.GetStub().DelState(accountID)
}

func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, fromId, toId string, amount float32) error {
	// @todo q solo pueda hacer esto el duenyo d la cuenta fuente

	// Get client org id and verify it matches peer org id.
	clientOrgID, err := getClientOrgID(ctx)
	if err != nil {
		return err
	}
	err = verifyClientOrgMatchesPeerOrg(clientOrgID)
	if err != nil {
		return err
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	fromAcc, err := s.ReadAccount(ctx, fromId)
	if err != nil {
		return errors.New("the source account doesn't exist")
	}

	toAcc, err := s.ReadAccount(ctx, toId)
	if err != nil {
		return errors.New("the destination account doesn't exist")
	}

	if fromAcc.Balance-amount < 0 {
		return errors.New("the source account does not have enough balance")
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

func (s *SmartContract) GetAllTxs(ctx contractapi.TransactionContextInterface, accountID string) ([]TxRecord, error) {

	// Get client org id and verify it matches peer org id.
	clientOrgID, err := getClientOrgID(ctx)
	if err != nil {
		return nil, err
	}
	err = verifyClientOrgMatchesPeerOrg(clientOrgID)
	if err != nil {
		return nil, err
	}

	// Get the transaction history result of an account
	historyIterator, err := ctx.GetStub().GetHistoryForKey(accountID)
	if err != nil {
		return nil, err
	}
	defer historyIterator.Close()

	var records []TxRecord
	for historyIterator.HasNext() {
		response, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}

		var account Account
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &account)
			if err != nil {
				return nil, err
			}
		} else {
			account = Account{
				ID: accountID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := TxRecord{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &account,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}
