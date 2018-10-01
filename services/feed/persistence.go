package feed

import (
	"encoding/json"

	"github.com/ppincak/rysen/pkg/persistence"
	"github.com/syndtr/goleveldb/leveldb"
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
func (persistence *Persistence) SaveFeed(metadata *Metadata) error {
	return persistence.Persist(metadata.Name, metadata)
}

// Get Feeds
func (persistence *Persistence) GetFeeds() ([]*Metadata, error) {
	feeds := make([]*Metadata, 0)
	err := persistence.Iterate(func(key []byte, value []byte) error {
		var feed *Metadata
		err := json.Unmarshal(value, &feed)
		if err != nil {
			return err
		}
		feeds = append(feeds, feed)
		return nil
	})
	return feeds, err
}
