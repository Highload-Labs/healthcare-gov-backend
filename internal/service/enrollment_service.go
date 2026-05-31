package service

import (
	"context"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

type EnrollmentService interface {
	EnrollPlan(ctx context.Context, input EnrollPlanInput) (*domain.Enrollment, error)
}

type EnrollmentServiceImpl struct {
	repo repository.EnrollmentRepository
}

func NewEnrollmentService(repo repository.EnrollmentRepository) EnrollmentService {
	return &EnrollmentServiceImpl{
		repo: repo,
	}
}

type EnrollPlanInput struct {
	UserID string
	PlanID string
}

func (s *EnrollmentServiceImpl) EnrollPlan(ctx context.Context, input EnrollPlanInput) (*domain.Enrollment, error) {
	effectiveDate := time.Now()
	endDate := effectiveDate.AddDate(1, 0, 0)

	newEnrollment := domain.Enrollment{
		UserID:        input.UserID,
		PlanID:        input.PlanID,
		EffectiveDate: effectiveDate,
		EndDate:       endDate,
	}

	id, err := s.repo.Create(ctx, newEnrollment)
	if err != nil {
		return nil, err
	}

	newEnrollment.ID = id

	return &newEnrollment, nil
}
