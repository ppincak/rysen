package schema

import (
	"encoding/json"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

type SchemaBackup struct {
	prefix string
	dbPath string
	db     *leveldb.DB
}

// Create Backup helper
func NewSchemaBackup(dbPath string, prefix *string) *SchemaBackup {
	if prefix == nil {
		_prefix := "schemas"
		prefix = &_prefix
	}

	return &SchemaBackup{
		prefix: *prefix,
		dbPath: dbPath,
	}
}

// Open connection
func (backup *SchemaBackup) Open() (err error) {
	db, err := leveldb.OpenFile(backup.dbPath, nil)
	if err != nil {
		return
	}
	backup.db = db
	return
}

// Close connection
func (backup *SchemaBackup) Close() error {
	if backup.db != nil {
		return backup.db.Close()
	}
	return nil
}

// Assemble the key
func (backup *SchemaBackup) assembleKey(key string) []byte {
	return []byte(strings.Join([]string{backup.prefix, "/", key}, ""))
}

// Persist schema
func (backup *SchemaBackup) SaveSchema(schema *ExchangeSchemaMetadata) (err error) {
	marshalled, err := json.Marshal(schema)
	if err != nil {
		return
	}
	err = backup.db.Put(backup.assembleKey(schema.Name), marshalled, nil)
	if err != nil {
		return
	}
	return nil
}

// Get all persisted schemas
func (backup *SchemaBackup) GetSchemas() ([]*ExchangeSchemaMetadata, error) {
	iterator := backup.db.NewIterator(nil, nil)
	defer iterator.Release()

	list := make([]*ExchangeSchemaMetadata, 0)
	for ok := iterator.Seek([]byte(backup.prefix)); ok; ok = iterator.Next() {
		var exchangeSchema *ExchangeSchemaMetadata
		err := json.Unmarshal(iterator.Value(), &exchangeSchema)
		if err != nil {
			return nil, err
		}
		list = append(list, exchangeSchema)
	}
	err := iterator.Error()
	if err != nil {
		return nil, err
	}
	return list, nil
}
