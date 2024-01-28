package main

import (
	"HardwareHunt/internal/domain"
	"HardwareHunt/internal/infrastructure/compday_ru"
	"HardwareHunt/internal/repository"
	"HardwareHunt/internal/service"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	//Настройка лога файлов
	setupLogger()
	// Инициализация конфигурационного файла
	if err := initConfig(); err != nil {
		log.Fatalf("Ошибка инициализации конфиг файла %s", err.Error())
		return
	}
	log.Println("Конфигурационный файл успешно инициализирован")

	// Загрузка переменных среды из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %s", err)
		return
	}
	log.Println("Файл .env успешно загружен")

	// Инициализация подключения к базе данных PostgresSQL
	db, err := initDatabaseConnection()
	if err != nil {
		log.Fatalf("Ошибка инициализации подключения к базе данных: %s", err.Error())
		return
	}
	log.Println("Подключение к базе данных успешно установлено")

	// Инициализация репозиториев, сервисов и обработчиков
	// ↑3)Работа с БД
	repos := repository.NewRepository(db)
	log.Println("Репозитории успешно инициализированы")

	// ↑2)Бизнес логика
	services := service.NewService(repos)
	log.Println("Бизнес-логика успешно инициализирована")

	processCompDayProcessors(services)
}

func processCompDayProcessors(services *service.Service) {
	// ↑1) Логика парсера данных
	url := "https://www.compday.ru/komplektuyuszie/protsessory/"
	scraperCompdayRu := compday_ru.NewScraper(url, &compday_ru.ProcessorScraper{})

	links, err := scraperCompdayRu.ScrapeLinks(url)
	if err != nil {
		log.Println("Ошибка парсера ссылок scraperCompdayRu")
		return
	}

	for i := 0; i < len(links); i += 100 {
		end := i + 100
		if end > len(links) {
			end = len(links)
		}

		log.Println(fmt.Sprintf("Начинаем парсить данные [%d/%d]. Всего %d", i, end, len(links)))
		iProcessors, err := scraperCompdayRu.ScrapeInfoFromLinks(links[i:end])
		if err != nil {
			log.Println(fmt.Sprintf("Ошибка парсинга данных [%d/%d]. Всего %d", i, end, len(links)))
			continue
		}
		log.Println(fmt.Sprintf("Данные успешно спарсены [%d/%d]. Всего %d", i, end, len(links)))

		var processors []domain.Processor
		for _, item := range iProcessors {
			processor, ok := item.(domain.Processor)
			if !ok {
				log.Println("Processors имеет неверный тип (!domain.Processor): ", item)
				return
			}
			processors = append(processors, processor)
		}

		// Пакетом вставляем новые записи
		numberRecordsAdded, err := services.CreateBatch(processors)
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Записи успешно вставлены. Количество: [%d]", numberRecordsAdded))
		}

		// Пакетом обновляем записи
		numberUpdatedRecords, err := services.UpdateBatch(processors)
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Записи успешно обновлены. Количество: [%d]", numberUpdatedRecords))
		}
	}
}

func initDatabaseConnection() (*sqlx.DB, error) {
	return repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
}

// Метод initConfig используется для инициализации конфигурационного файла
func initConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.SetConfigType("yaml") // Добавьте эту строку, если файл имеет расширение .yaml
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации: %s", err)
		return err
	}
	return nil
}

// setupLogger создает логгер и настраивает форматирование журнала с использованием log
func setupLogger() {
	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл журнала: %s", err.Error())
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.JSONFormatter{})
}
