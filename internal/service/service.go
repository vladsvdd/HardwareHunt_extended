package service

import (
	"HardwareHunt/internal/domain"
	"HardwareHunt/internal/repository"
)

type PcProcessor interface {
	Create(processor domain.Processor) (int, error)
	CreateBatch(processors []domain.Processor) (int, error)
	UpdateBatch(processors []domain.Processor) (int, error)
}

type Service struct {
	PcProcessor
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		PcProcessor: NewPcProcessorService(repos.PcProcessor),
	}
}
