package service

import (
	"HardwareHunt/internal/domain"
	"HardwareHunt/internal/repository"
)

type PcProcessorService struct {
	repo repository.PcProcessor
}

func NewPcProcessorService(repo repository.PcProcessor) *PcProcessorService {
	return &PcProcessorService{repo: repo}
}

func (s *PcProcessorService) Create(processor domain.Processor) (int, error) {
	return s.repo.Create(processor)
}

func (s *PcProcessorService) CreateBatch(processors []domain.Processor) (int, error) {
	return s.repo.CreateBatch(processors)
}

func (s *PcProcessorService) UpdateBatch(processors []domain.Processor) (int, error) {
	return s.repo.UpdateBatch(processors)
}
