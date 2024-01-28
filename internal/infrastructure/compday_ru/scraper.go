package compday_ru

import "github.com/gocolly/colly"

// Scraper - основной тип, который реализует Scrappable
type Scraper struct {
	url         string
	IScrappable Scrappable
}

// Scrappable - общий интерфейс для данных, которые мы хотим парсить
type Scrappable interface {
	ScrapeInfo(c *colly.Collector, link string) (interface{}, []error)
}

func NewScraper(url string, scrap Scrappable) *Scraper {
	return &Scraper{
		url:         url,
		IScrappable: scrap,
	}
}
