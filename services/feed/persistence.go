package feed

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
	"rysen/pkg/persistence"
)

type Persistence struct {
	*persistence.Persistence
}

// Create new Feed backup
func NewPersistence(db *leveldb.DB, prefix *string) *Persistence {
	if prefix == nil {
		_prefix := "feeds"
		prefix = &_prefix
	}

	return &Persistence{
		Persistence: persistence.NewPersistence(db, *prefix),
	}
}

// Save feed
func (persistence *Persistence) SaveFeed(model *Model) error {
	return persistence.Persist(model.Name, model)
}

// Get Feeds
func (persistence *Persistence) GetFeeds() ([]*Model, error) {
	feeds := make([]*Model, 0)
	err := persistence.Iterate(func(key []byte, value []byte) error {
		var feed *Model
		err := json.Unmarshal(value, &feed)
		if err != nil {
			return err
		}
		feeds = append(feeds, feed)
		return nil
	})
	return feeds, err
}
