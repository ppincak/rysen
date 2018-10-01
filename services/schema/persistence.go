package schema

import (
	"encoding/json"

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

// Persist schema
func (persistence *Persistence) SaveSchema(schema *ExchangeSchemaMetadata) (err error) {
	return persistence.Persist(schema.Name, schema)
}

// Get all persisted schemas
func (persistence *Persistence) GetSchemas() ([]*ExchangeSchemaMetadata, error) {
	list := make([]*ExchangeSchemaMetadata, 0)
	err := persistence.Iterate(func(key []byte, value []byte) error {
		var exchangeSchema *ExchangeSchemaMetadata
		err := json.Unmarshal(value, &exchangeSchema)
		if err != nil {
			return err
		}
		list = append(list, exchangeSchema)
		return nil
	})
	return list, err
}
