package service

import (
	"bwg/internal/models"
	"log/slog"
)

type repository interface {
	NewTicker(ticker string) error
	GetPriceDifference(info models.TickerInfo) (models.TicketDifference, error)
}

type Service struct {
	repository
}

func New(repository repository) *Service {
	return &Service{repository}
}

func (s *Service) CreateNewTicker(ticker string) error {
	err := s.repository.NewTicker(ticker)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetTickerInfo(info models.TickerInfo) (models.TicketDifference, error) {
	slog.Info("Service get:", info)
	return s.repository.GetPriceDifference(info)
}
