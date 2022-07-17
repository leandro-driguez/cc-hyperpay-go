package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lllrdgz/cc-hyperpay-go/hyperpay-transfer/chaincode"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type HyperPayContract struct {
	c *gateway.Contract
}

func NewHyperPayContract() (*HyperPayContract, error) {
	const (
		channelId = "mychannel"
		identity  = "User1@org1.example.com"
	)
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		return nil, err
	}
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return nil, err
	}
	if !wallet.Exists(identity) {
		if err := populateWallet(wallet); err != nil {
			return nil, err
		}
	}
	ccpPath := filepath.Join(
		"ccp.yaml",
	)
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, identity),
	)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelId)
	if err != nil {
		return nil, err
	}
	return &HyperPayContract{c: network.GetContract("mycc")}, nil
}

func (contract *HyperPayContract) Init() error {
	_, err := contract.c.SubmitTransaction("InitLedger")
	if err != nil {
		return err
	}
	return nil
}

func (contract *HyperPayContract) Read(id string) (*chaincode.Account, error) {
	result, err := contract.c.EvaluateTransaction("ReadAccount", id)
	if err != nil {
		return nil, err
	}
	var account chaincode.Account
	err = json.Unmarshal(result, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (contract HyperPayContract) Transfer(fromId, toId string, amount float32) error {
	_, err := contract.c.SubmitTransaction("Transfer", fromId, toId, fmt.Sprint(amount))
	if err != nil {
		return err
	}
	return nil
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}
