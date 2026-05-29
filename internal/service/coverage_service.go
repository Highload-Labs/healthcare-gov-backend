package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

type CoverageService interface {
	GetCoverageByZipcode(ctx context.Context, zipcode string) (*domain.Coverage, error)
}

var ErrCoverageNotFound = errors.New("coverage not found")

type CoverageServiceImpl struct {
	repo repository.CoverageRepository
}

func NewCoverageService(repo repository.CoverageRepository) CoverageService {
	return &CoverageServiceImpl{
		repo: repo,
	}
}

func (s *CoverageServiceImpl) GetCoverageByZipcode(ctx context.Context, zipcode string) (*domain.Coverage, error) {
	coverage, err := s.repo.FindByZipcode(ctx, zipcode)
	if err != nil {
		if errors.Is(err, repository.ErrCoverageNotFound) {
			return nil, ErrCoverageNotFound
		}

		return nil, err
	}

	return coverage, nil
}
