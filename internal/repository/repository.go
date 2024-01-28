package repository

import (
	"HardwareHunt/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PcProcessor interface {
	Create(processor domain.Processor) (int, error)
	CreateBatch(processors []domain.Processor) (int, error)
	UpdateBatch(processors []domain.Processor) (int, error)
}

type Repository struct {
	PcProcessor
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PcProcessor: NewPcProcessorPostgres(db),
	}
}
