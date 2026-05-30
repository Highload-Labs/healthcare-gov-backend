package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

type PlanService interface {
	GetByZipcode(ctx context.Context, zipcode string, pagination *shared.Pagination) (
		[]domain.Plan,
		*shared.Metadata,
		error,
	)
	GetById(ctx context.Context, id string) (*domain.Plan, error)
}

var ErrPlanNotFound = errors.New("plan not found")

type PlanServiceImpl struct {
	repo repository.PlanRepository

	coverageService CoverageService
}

func NewPlanService(repo repository.PlanRepository, coverageService CoverageService) PlanService {
	return &PlanServiceImpl{
		repo:            repo,
		coverageService: coverageService,
	}
}

func (s *PlanServiceImpl) GetByZipcode(
	ctx context.Context,
	zipcode string,
	pagination *shared.Pagination,
) ([]domain.Plan, *shared.Metadata, error) {
	coverage, err := s.coverageService.GetCoverageByZipcode(ctx, zipcode)
	if err != nil {
		return nil, nil, err
	}

	offset := pagination.CalculateOffset()
	metadata := &shared.Metadata{
		CurrentPage: pagination.PageNumber,
		Limit:       pagination.Limit,
	}

	state := coverage.State
	totalData, err := s.repo.CountByState(ctx, state)
	if err != nil {
		return nil, nil, err
	}

	if totalData == 0 {
		return nil, nil, ErrPlanNotFound
	}

	plans, err := s.repo.FindByState(ctx, state, pagination.Limit, offset)
	if err != nil {
		if errors.Is(err, repository.ErrPlanNotFound) {
			return nil, nil, ErrPlanNotFound
		}

		return nil, nil, err
	}

	totalPages := pagination.CalculateTotalPages(totalData)
	metadata.TotalData = totalData
	metadata.TotalPages = totalPages

	return plans, metadata, nil
}

func (s *PlanServiceImpl) GetById(ctx context.Context, id string) (*domain.Plan, error) {
	plan, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPlanNotFound) {
			return nil, ErrPlanNotFound
		}
	}

	return plan, nil
}
