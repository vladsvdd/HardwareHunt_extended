Вот пример структуры проекта для вашего веб-скрапинга с использованием чистой архитектуры:

```plaintext
books-scraper/
|-- cmd/
|   `-- scraper/
|       `-- main.go
|-- internal/
|   |-- domain/
|   |   `-- book.go
|   |-- infrastructure/
|   |   `-- scraper_collector.go
|   `-- interfaces/
|       `-- csv_writer.go
|-- .gitignore
|-- go.mod
|-- go.sum
|-- README.md
```

Объяснение:

- `cmd/scraper/main.go`: Основной файл для запуска приложения.
- `internal/domain/`: Здесь располагаются бизнес-модели (например, `Book`).
- `internal/infrastructure/`: Реализации инфраструктурных компонентов (например, `ScraperCollector`).
- `internal/interfaces/`: Компоненты, связанные с интерфейсами (например, `CSVWriter`).
- `.gitignore`: Игнорируемые файлы для системы контроля версий Git.
- `go.mod` и `go.sum`: Файлы модуля Go, описывающие зависимости проекта.
- `README.md`: Документация о проекте.

Вам также нужно создать соответствующие файлы для `book.go`, `scraper_collector.go` и `csv_writer.go` в соответствующих каталогах. Эти файлы будут содержать код, описывающий бизнес-логику, инфраструктурные компоненты и интерфейсные компоненты соответственно.


Проверяем подключение к selenium
```
curl http://localhost:4444/wd/hub/status
```
