package persistence

import (
	"encoding/json"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

type Persistence struct {
	prefix string
	db     *leveldb.DB
}

// Create new backup
func NewPersistence(db *leveldb.DB, prefix string) *Persistence {
	return &Persistence{
		db:     db,
		prefix: prefix,
	}
}

// Level Db
func (persistence *Persistence) Db() *leveldb.DB {
	return persistence.db
}

// Prefix
func (persistence *Persistence) Prefix() string {
	return persistence.prefix
}

// Assemble the key
func (persistence *Persistence) AssembleKey(key string) []byte {
	return []byte(strings.Join([]string{persistence.prefix, "/", key}, ""))
}

// Persist
func (persistence *Persistence) Persist(key string, value interface{}) (err error) {
	marshalled, err := json.Marshal(value)
	if err != nil {
		return
	}
	err = persistence.Db().Put(persistence.AssembleKey(key), marshalled, nil)
	if err != nil {
		return
	}
	return nil
}

// Delete
func (persistence *Persistence) Delete(key string) (err error) {
	return persistence.Db().Delete(persistence.AssembleKey(key), nil)
}

// Get persited value
func (persistence *Persistence) Get(key string, marshalled interface{}) (err error) {
	value, err := persistence.Db().Get(persistence.AssembleKey(key), nil)
	if err != nil {
		return
	}
	return json.Unmarshal(value, marshalled)
}

// Iteration handler function
type IterateHandler func(key []byte, value []byte) error

// Get all persisted schemas
func (persistence *Persistence) Iterate(handler IterateHandler) (err error) {
	iterator := persistence.Db().NewIterator(nil, nil)
	defer iterator.Release()

	for ok := iterator.Seek([]byte(persistence.Prefix())); ok; ok = iterator.Next() {
		err = handler(iterator.Key(), iterator.Value())
		if err != nil {
			return
		}
	}
	err = iterator.Error()
	if err != nil {
		return
	}
	return
}
