package scraper

type ScraperMetadata struct {
	Topic    string   `json:"topic"`
	Symbols  []string `json:"symbols"`
	Interval int64    `json:"interval"`
}
