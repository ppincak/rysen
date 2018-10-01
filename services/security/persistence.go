package security

import (
	"encoding/json"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/ppincak/rysen/pkg/persistence"
)

type Persistence struct {
	*persistence.Persistence
}

// Create Backup helper
func NewPersistence(db *leveldb.DB, prefix *string) *Persistence {
	if prefix == nil {
		_prefix := "schemas"
		prefix = &_prefix
	}

	return &Persistence{
		Persistence: persistence.NewPersistence(db, *prefix),
	}
}

// TODO: add encryption
// Persist schema
func (persistence *Persistence) SaveAccount(account *Account) (err error) {
	return persistence.Persist(strings.Join([]string{account.Exchange, account.Name}, "/"), account)
}

// TODO: add decryption
// Get all persisted schemas
func (persistence *Persistence) GetAccounts() ([]*Account, error) {
	list := make([]*Account, 0)
	err := persistence.Iterate(func(key []byte, value []byte) error {
		var account *Account
		err := json.Unmarshal(value, &account)
		if err != nil {
			return err
		}
		list = append(list, account)
		return nil
	})
	return list, err
}
