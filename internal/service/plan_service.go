package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

type PlanService interface {
	GetByZipcode(ctx context.Context, zipcode string) ([]domain.Plan, error)
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

func (s *PlanServiceImpl) GetByZipcode(ctx context.Context, zipcode string) ([]domain.Plan, error) {
	coverage, err := s.coverageService.GetCoverageByZipcode(ctx, zipcode)
	if err != nil {
		return nil, err
	}

	state := coverage.State
	plans, err := s.repo.FindByState(ctx, state)
	if err != nil {
		if errors.Is(err, repository.ErrPlanNotFound) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	if len(plans) == 0 {
		return nil, ErrPlanNotFound
	}

	return plans, nil
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
