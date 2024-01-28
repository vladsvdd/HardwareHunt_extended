// Package compday_ru links_scraper.go // Файл с методами для парсинга ссылок
package compday_ru

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

func (s *Scraper) ScrapeLinks(url string) (links []string, err error) {
	c := colly.NewCollector()
	// Извлечение блока с процессорами
	c.OnHTML("#catItems", func(element *colly.HTMLElement) {
		// Извлечение ссылки на процессоры
		links = element.ChildAttrs("a.name", "href")
	})

	err = c.Visit(url)
	return links, err
}

// ScrapeInfoFromLinks парсит информацию из списка ссылок
func (s *Scraper) ScrapeInfoFromLinks(links []string) ([]interface{}, error) {
	var data []interface{}

	var wg sync.WaitGroup
	source := rand.NewSource(time.Now().UnixNano())
	rateLimit := time.Duration(rand.New(source).Intn(3)+1) * time.Second
	// Создайте канал для управления скоростью
	rateLimitChannel := time.Tick(rateLimit)

	c := colly.NewCollector()
	for _, link := range links {
		// Ожидайте, чтобы соблюсти ограничение скорости
		<-rateLimitChannel

		// Увеличьте счетчик WaitGroup перед каждым запросом
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			log.Printf("Начинается обработка ссылки: %s", link)

			// Реализуйте ваш код для запроса по ссылке и обработки данных
			dataItem, err := s.IScrappable.ScrapeInfo(c, link)
			if err != nil {
				log.Printf("Ошибки при обработке ссылки %s:", link)
				for _, e := range err {
					log.Print("  - ", e.Error())
				}
			} else {
				data = append(data, dataItem)
				log.Printf("Обработка ссылки %s завершена успешно", link)
			}
		}(link)
	}

	// Ожидайте завершения всех запросов
	wg.Wait()

	return data, nil
}
