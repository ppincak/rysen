package scraper

type Metadata struct {
	Topic      string   `json:"topic"`
	Symbols    []string `json:"symbols"`
	Interval   int64    `json:"interval"`
	ScrapeFunc string   `json:"scrapeFunction"`
}
