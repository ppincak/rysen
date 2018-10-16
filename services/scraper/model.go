package scraper

import (
	"github.com/ppincak/rysen/pkg/collections"
)

// Scraper model
type Model struct {
	Topic      string   `json:"topic"`
	Symbols    []string `json:"symbols"`
	Interval   int64    `json:"interval"`
	ScrapeFunc string   `json:"scrapeFunction"`
}

var _ collections.Comparable = (*Model)(nil)

// Equals func
func (model *Model) Equals(value interface{}) bool {
	assertion, ok := value.(*Model)
	if !ok {
		return false
	}
	if assertion.Topic != model.Topic {
		return false
	}
	return true
}
