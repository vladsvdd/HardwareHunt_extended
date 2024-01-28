// Package compday_ru Файл с методами для парсинга информации о процессорах
package compday_ru

import (
	"HardwareHunt/internal/domain"
	"errors"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

// ProcessorScraper - тип-обертка для domain. Processor
type ProcessorScraper struct {
	*domain.Processor
}

// ProcessorScraper - это тип данных, который мы хотим убедиться, что реализует интерфейс Scrappable
var _ Scrappable = &ProcessorScraper{}

// ScrapeInfo реализует метод ScrapeInfo интерфейса Scrappable для Processor
// Парсинга информации о процессорах
func (p *ProcessorScraper) ScrapeInfo(c *colly.Collector, link string) (interface{}, []error) {
	startUrl := "https://www.compday.ru"
	var processor domain.Processor
	var errLogs []error
	// Парсинг таблицы с информацией о процессоре
	c.OnHTML("table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			key := el.ChildText("td:nth-child(1)")
			value := el.ChildText("td:nth-child(2)")
			localErr := p.makeObjectProcessor(&processor, key, value)
			if localErr != nil {
				errLogs = append(errLogs, localErr)
			}
		})
	})

	// Парсинг новой цены
	c.OnHTML(".actual.cash.price", func(e *colly.HTMLElement) {
		localErr := p.makeObjectProcessor(&processor, ".actual.cash.price", e.Text)
		if localErr != nil {
			errLogs = append(errLogs, localErr)
		}
	})

	// Парсинг старой цены
	c.OnHTML(".old.price", func(e *colly.HTMLElement) {
		localErr := p.makeObjectProcessor(&processor, ".old.price", e.Text)
		if localErr != nil {
			errLogs = append(errLogs, localErr)
		}
	})

	// Обработчик для поиска слова "Товар доступен"
	c.OnHTML("p:contains('Товар доступен') strong", func(e *colly.HTMLElement) {
		localErr := p.makeObjectProcessor(&processor, "Товар доступен", e.Text)
		if localErr != nil {
			errLogs = append(errLogs, localErr)
		}
	})

	// Обработчик города
	c.OnHTML("a.city", func(e *colly.HTMLElement) {
		localErr := p.makeObjectProcessor(&processor, "city", e.Text)
		if localErr != nil {
			errLogs = append(errLogs, localErr)
		}
	})

	// Обработчик улицы
	c.OnHTML("a.street", func(e *colly.HTMLElement) {
		localErr := p.makeObjectProcessor(&processor, "street", e.Text)
		if localErr != nil {
			errLogs = append(errLogs, localErr)
		}
	})

	err := p.makeObjectProcessor(&processor, "link", startUrl+link)
	if err != nil {
		errLogs = append(errLogs, err)
	}

	err = c.Visit(startUrl + link)
	if err != nil {
		errLogs = append(errLogs, errors.New("ошибка посещения страницы. Метод ScrapeInfoFromLinks. "+err.Error()))
	}

	return processor, errLogs
}

func (p *ProcessorScraper) makeObjectProcessor(processor *domain.Processor, key, value string) error {
	switch key {
	case "Линейка":
		processor.Line = value
	case "Модель":
		processor.Model = value
	case "Ядро":
		processor.Core = value
	case "Количество ядер":
		numCores, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return errors.New("ошибка конвертации \"Количество ядер\"." + err.Error())
		}
		processor.NumCores = int(numCores)
	case "Производитель":
		processor.Manufacturer = value
	case "Код производителя":
		processor.ManufacturerCode = value
	case "Socket":
		processor.Socket = value
	case "Типичное тепловыделение":
		thermalPower, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return errors.New("ошибка конвертации \"Типичное тепловыделение\"." + err.Error())
		}
		processor.ThermalPower = int(thermalPower)
	case "Сайт производителя":
		processor.ManufacturerWebsite = value
	case "Технологический процесс":
		technologyProcess, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return errors.New("ошибка конвертации \"Технологический процесс\"." + err.Error())
		}
		processor.TechnologyProcess = int(technologyProcess)
	case "Технологии":
		processor.Technologies = value
	case "Тип поставки":
		processor.DeliveryType = value
	case "Объем кэша L1":
		l1CacheVolume, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("ошибка конвертации \"Объем кэша L1\"." + err.Error())
		}
		processor.L1CacheVolume = l1CacheVolume
	case "Объем кэша L2":
		l2CacheVolume, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("ошибка конвертации \"Объем кэша L2\"." + err.Error())
		}
		processor.L2CacheVolume = l2CacheVolume
	case "Объем кэша L3":
		l3CacheVolume, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("ошибка конвертации \"Объем кэша L3\"." + err.Error())
		}
		processor.L3CacheVolume = l3CacheVolume
	case "Интегрированное графическое ядро":
		if "да" == strings.ToLower(value) {
			processor.IntegratedGraphics = true
		} else {
			processor.IntegratedGraphics = false
		}
	case "Видеопроцессор":
		processor.VideoProcessor = value
	case ".actual.cash.price":
		newPrice := strings.ReplaceAll(value, " ", "")
		parseNewPrice, err := strconv.ParseFloat(newPrice, 64)
		if err != nil {
			err = errors.New("ошибка преобразования новой цены. Метод ScrapeInfoFromLinks. " + err.Error())
		}
		processor.NewPrice = parseNewPrice
	case ".old.price":
		oldPrice := strings.ReplaceAll(value, " ", "")
		parseOldPrice, err := strconv.ParseFloat(oldPrice, 64)
		if err != nil {
			err = errors.New("ошибка преобразования старой цены. Метод ScrapeInfoFromLinks. " + err.Error())
		}
		processor.OldPrice = parseOldPrice
	case "link":
		processor.Link = value
	case "Товар доступен":
		processor.AvailabilityDate = strings.TrimSpace(value)
	case "city":
		processor.City = strings.TrimSpace(value)
	case "street":
		processor.Street = strings.TrimSpace(value)
	}

	return nil
}
