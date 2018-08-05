package services

import (
	"github.com/ppincak/rysen/pkg/scrape"
)

type ScraperService struct {
	scrapers map[string]*scrape.Scraper
}

func NewScraperService() {

}
