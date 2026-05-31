package service

import (
	"context"
	"errors"
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

var InvalidUserOrPlan = errors.New("invalid user or plan")
var ErrUserHasActivePlan = errors.New("user already has active plan")

func (s *EnrollmentServiceImpl) EnrollPlan(ctx context.Context, input EnrollPlanInput) (*domain.Enrollment, error) {
	enrollment, err := s.repo.FindActiveEnrollmentByUserID(ctx, input.UserID)
	if err != nil {
		if !errors.Is(err, repository.ErrNoActivePlan) {
			return nil, err
		}
	}

	if enrollment != nil {
		return nil, ErrUserHasActivePlan
	}

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
		if errors.Is(err, repository.InvalidUserOrPlan) {
			return nil, InvalidUserOrPlan
		}

		return nil, err
	}

	newEnrollment.ID = id

	return &newEnrollment, nil
}
