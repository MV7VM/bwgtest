package service

type repository interface {
	NewTicker(ticker string) error
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
