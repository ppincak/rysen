package services

import (
	"github.com/ppincak/rysen/pkg/scrape"
)

type ScraperMetadata struct {
	Topic      string                `json:"topic"`
	Symbols    []string              `json:"symbols"`
	Interval   int64                 `json:"interval"`
	ScrapeFunc scrape.CallerFuncType `json:"scrapeFunc"`
}
