package repository

import (
	"HardwareHunt/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type PcProcessorPostgres struct {
	db *sqlx.DB
}

func NewPcProcessorPostgres(db *sqlx.DB) *PcProcessorPostgres {
	return &PcProcessorPostgres{db: db}
}

func (r *PcProcessorPostgres) Create(processor domain.Processor) (int, error) {
	//используем транзакцию
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	// Генерируем строку запроса с именованными параметрами
	createProcessQuery := fmt.Sprintf(`
		INSERT INTO %s (
			link, line, model, core, num_cores, manufacturer, manufacturer_code, socket,
			thermal_power, manufacturer_website, technology_process, technologies, delivery_type,
			l1_cache_volume, l2_cache_volume, l3_cache_volume, integrated_graphics, video_processor,
			new_price, old_price, availability_date, city, street
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23
		) RETURNING id
	`, pcProcessorTable)

	// Выполняем запрос с именованными аргументами
	row := transaction.QueryRow(createProcessQuery, processor.Link, processor.Line, processor.Model, processor.Core,
		processor.NumCores, processor.Manufacturer, processor.ManufacturerCode, processor.Socket, processor.ThermalPower,
		processor.ManufacturerWebsite, processor.TechnologyProcess, processor.Technologies, processor.DeliveryType,
		processor.L1CacheVolume, processor.L2CacheVolume, processor.L3CacheVolume, processor.IntegratedGraphics,
		processor.VideoProcessor, processor.NewPrice, processor.OldPrice, processor.AvailabilityDate, processor.City,
		processor.Street)
	err = row.Scan(&id)
	if err != nil {
		transaction.Rollback()
		return 0, errors.New(fmt.Sprintf("Процессор %v уже существует в БД. %v", processor.Link, err.Error()))
	}

	return id, transaction.Commit()
}

func (r *PcProcessorPostgres) CreateBatch(processors []domain.Processor) (int, error) {
	var numberRecordsAdded int
	if 0 == len(processors) {
		return numberRecordsAdded, nil
	}

	// Начать транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		return numberRecordsAdded, err
	}
	defer tx.Rollback()

	// Собрать все значения link
	var links []string
	for _, p := range processors {
		links = append(links, p.Link)
	}

	// Выполнить запрос для проверки наличия записей с link
	var existingLinks []string
	err = tx.Select(&existingLinks, fmt.Sprintf("SELECT link FROM %s WHERE link = ANY ($1)", pcProcessorTable), pq.Array(links))
	if err != nil {
		return numberRecordsAdded, errors.New("Ошибка при выполнении запроса проверки наличия записей с link. " + err.Error())
	}

	// Проверить, есть ли уже существующие записи
	existingSet := make(map[string]struct{})
	for _, existingLink := range existingLinks {
		existingSet[existingLink] = struct{}{}
	}

	// Создать новый срез для хранения уникальных links
	var uniqueProcessors []domain.Processor

	// Перебрать все links и добавить только уникальные в новый срез
	for _, processor := range processors {
		if _, exists := existingSet[processor.Link]; !exists {
			uniqueProcessors = append(uniqueProcessors, processor)
		}
	}

	// Теперь processors содержит только уникальные записи
	processors = uniqueProcessors

	if processors == nil {
		tx.Rollback()
		return numberRecordsAdded, nil
	}

	// Подготовить запрос с именованными параметрами
	query := `
		INSERT INTO pc_processor (
			link, line, model, core, num_cores, manufacturer, manufacturer_code, socket,
			thermal_power, manufacturer_website, technology_process, technologies, delivery_type,
			l1_cache_volume, l2_cache_volume, l3_cache_volume, integrated_graphics, video_processor,
			new_price, old_price, availability_date, city, street
		) VALUES (
			:link, :line, :model, :core, :num_cores, :manufacturer, :manufacturer_code, :socket,
			:thermal_power, :manufacturer_website, :technology_process, :technologies, :delivery_type,
			:l1_cache_volume, :l2_cache_volume, :l3_cache_volume, :integrated_graphics, :video_processor,
			:new_price, :old_price, :availability_date, :city, :street
		)
	`

	// Использовать NamedExec для вставки всех записей за один раз
	_, err = tx.NamedExec(query, processors)
	if err != nil {
		return numberRecordsAdded, err
	}

	numberRecordsAdded = len(processors)

	// Если все успешно, подтвердить транзакцию
	return numberRecordsAdded, tx.Commit()
}

func (r *PcProcessorPostgres) UpdateBatch(processors []domain.Processor) (int, error) {
	var numberUpdatedRecords int
	if 0 == len(processors) {
		return numberUpdatedRecords, nil
	}

	// Начать транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		return numberUpdatedRecords, err
	}
	defer tx.Rollback()

	// existingLinks - это срез ссылок, которые вы хотите проверить
	var existingLinks []string

	// Подготовить запрос для проверки существования записей
	query := fmt.Sprintf(`
		SELECT link 
		FROM %s 
		WHERE `,
		pcProcessorTable)

	// Подготовить значения для вставки в запрос
	var values []interface{}
	var parameters []string
	var counter int = 0
	for _, processor := range processors {
		parameters = append(parameters, fmt.Sprintf(`
			(link = $%d 
			AND (
				line, model, core, num_cores, manufacturer,
				manufacturer_code, socket, thermal_power, manufacturer_website, technology_process,
				technologies, delivery_type, l1_cache_volume, l2_cache_volume, l3_cache_volume,
				integrated_graphics, video_processor, new_price, old_price, availability_date,
				city, street
			) IS DISTINCT FROM
			ROW($%d, $%d, $%d, $%d, $%d,
				$%d, $%d, $%d, $%d, $%d,
				$%d, $%d, $%d, $%d, $%d,
				$%d, $%d, $%d, $%d, $%d,
				$%d, $%d))
		`, counterPlusOne(&counter),
			counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter),
			counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter),
			counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter),
			counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter), counterPlusOne(&counter),
			counterPlusOne(&counter), counterPlusOne(&counter)))

		values = append(values,
			processor.Link,
			processor.Line, processor.Model, processor.Core, processor.NumCores, processor.Manufacturer,
			processor.ManufacturerCode, processor.Socket, processor.ThermalPower, processor.ManufacturerWebsite, processor.TechnologyProcess,
			processor.Technologies, processor.DeliveryType, processor.L1CacheVolume, processor.L2CacheVolume, processor.L3CacheVolume,
			processor.IntegratedGraphics, processor.VideoProcessor, processor.NewPrice, processor.OldPrice, processor.AvailabilityDate,
			processor.City, processor.Street)
	}

	query += strings.Join(parameters, " OR ")

	// Выполнить запрос
	err = tx.Select(&existingLinks, query, values...)

	if err != nil {
		return numberUpdatedRecords, errors.New("Ошибка при выполнении запроса проверки наличия записей с link. " + err.Error())
	}

	// existingLinks теперь содержит только те ссылки, где значения отличаются

	// Проверить, есть ли уже существующие записи
	existingSet := make(map[string]struct{})
	for _, existingLink := range existingLinks {
		existingSet[existingLink] = struct{}{}
	}

	// Создать новый срез для хранения уникальных links
	var existsProcessors []domain.Processor

	// Перебрать все links и добавить только уникальные в новый срез
	for _, processor := range processors {
		if _, exists := existingSet[processor.Link]; exists {
			existsProcessors = append(existsProcessors, processor)
		}
	}

	// Теперь processors содержит только уникальные записи
	processors = existsProcessors

	// Пройти по каждому процессору и обновить соответствующую запись
	for _, processor := range processors {
		// Подготовить запрос для обновления записи
		queryUpdate := `
        UPDATE pc_processor
        SET
            line = :line,
            model = :model,
            core = :core,
            num_cores = :num_cores,
            manufacturer = :manufacturer,
            manufacturer_code = :manufacturer_code,
            socket = :socket,
            thermal_power = :thermal_power,
            manufacturer_website = :manufacturer_website,
            technology_process = :technology_process,
            technologies = :technologies,
            delivery_type = :delivery_type,
            l1_cache_volume = :l1_cache_volume,
            l2_cache_volume = :l2_cache_volume,
            l3_cache_volume = :l3_cache_volume,
            integrated_graphics = :integrated_graphics,
            video_processor = :video_processor,
            new_price = :new_price,
            old_price = :old_price,
            availability_date = :availability_date,
            city = :city,
            street = :street
        WHERE link = :link
    `

		// Использовать NamedExec для обновления записи
		_, err = tx.NamedExec(queryUpdate, processor)
		if err != nil {
			return numberUpdatedRecords, err
		}
	}

	numberUpdatedRecords = len(processors)

	// Если все успешно, подтвердить транзакцию
	return numberUpdatedRecords, tx.Commit()
}

func counterPlusOne(num *int) int {
	*num += 1
	return *num
}
