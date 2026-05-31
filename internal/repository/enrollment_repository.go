package repository

import (
	"context"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment domain.Enrollment) (string, error)
}

type EnrollmentRepositoryImpl struct {
	pg *infra.Postgresql
}

func NewEnrollmentRepository(pg *infra.Postgresql) EnrollmentRepository {
	return &EnrollmentRepositoryImpl{pg}
}

func (r *EnrollmentRepositoryImpl) Create(ctx context.Context, enrollment domain.Enrollment) (string, error) {
	err := r.pg.Db.QueryRowContext(ctx, "INSERT INTO enrollments (user_id, plan_id, effective_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id", enrollment.UserID, enrollment.PlanID, enrollment.EffectiveDate, enrollment.EndDate).Scan(&enrollment.ID)
	if err != nil {
		return "", err
	}

	return enrollment.ID, nil
}
