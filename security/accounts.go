package security

import (
	"io/ioutil"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/pkg/encryption"
)

type Accounts struct {
	ExchangeAccounts map[string][]*Account
	Accounts         map[string]*Account
}

// Load accounts.json file
func LoadAndCreateAccounts(url string, decryptionKey []byte) (*Accounts, error) {
	bytes, err := ioutil.ReadFile(url)
	if err != nil {
		return nil, err
	}
	return CreateAccounts(bytes, decryptionKey)
}

// Create schema from json
func CreateAccounts(jsonSchema []byte, decryptionKey []byte) (*Accounts, error) {
	decoded, err := encryption.DecryptAES(jsonSchema, decryptionKey)
	if err != nil {
		return nil, err
	}

	var exchangeAccounts map[string][]*Account
	err = api.UnmarshallAs(decoded, &exchangeAccounts)
	if err != nil {
		return nil, err
	}

	accountsByNames := make(map[string]*Account)
	for _, accounts := range exchangeAccounts {
		for _, account := range accounts {
			accountsByNames[account.Name] = account
		}
	}

	return &Accounts{
		ExchangeAccounts: exchangeAccounts,
		Accounts:         accountsByNames,
	}, nil
}
