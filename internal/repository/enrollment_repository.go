package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/lib/pq"
)

type EnrollmentRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	LockByUserID(ctx context.Context, tx *sql.Tx, userID string) error
	Create(ctx context.Context, enrollment domain.Enrollment) (string, error)
	CreateTx(ctx context.Context, tx *sql.Tx, enrollment domain.Enrollment) (string, error)
	FindActiveEnrollmentByUserID(ctx context.Context, userID string) (*domain.Enrollment, error)
	FindActiveEnrollmentByUserIDTx(ctx context.Context, tx *sql.Tx, userID string, now time.Time) (
		*domain.Enrollment,
		error,
	)
}

type EnrollmentRepositoryImpl struct {
	pg *infra.Postgresql
}

var InvalidUserOrPlan = errors.New("invalid user or plan")
var ErrNoActivePlan = errors.New("no active plan")

func NewEnrollmentRepository(pg *infra.Postgresql) EnrollmentRepository {
	return &EnrollmentRepositoryImpl{pg}
}

func (r *EnrollmentRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.pg.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *EnrollmentRepositoryImpl) LockByUserID(ctx context.Context, tx *sql.Tx, userID string) error {
	_, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtext($1))", userID)
	if err != nil {
		return err
	}

	var now time.Time
	_ = tx.QueryRowContext(ctx, "SELECT now()").Scan(&now)
	slog.Info("now time based on TX", "time", now)

	return nil
}

func (r *EnrollmentRepositoryImpl) Create(ctx context.Context, enrollment domain.Enrollment) (
	string,
	error,
) {
	err := r.pg.Db.QueryRowContext(
		ctx,
		"INSERT INTO enrollments (user_id, plan_id, effective_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id",
		enrollment.UserID,
		enrollment.PlanID,
		enrollment.EffectiveDate,
		enrollment.EndDate,
	).Scan(&enrollment.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return "", InvalidUserOrPlan
			}
		}
		return "", err
	}

	return enrollment.ID, nil
}

func (r *EnrollmentRepositoryImpl) CreateTx(ctx context.Context, tx *sql.Tx, enrollment domain.Enrollment) (
	string,
	error,
) {
	err := tx.QueryRowContext(
		ctx,
		"INSERT INTO enrollments (user_id, plan_id, effective_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id",
		enrollment.UserID,
		enrollment.PlanID,
		enrollment.EffectiveDate,
		enrollment.EndDate,
	).Scan(&enrollment.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return "", InvalidUserOrPlan
			}
		}
		return "", err
	}

	return enrollment.ID, nil
}

func (r *EnrollmentRepositoryImpl) FindActiveEnrollmentByUserID(ctx context.Context, userID string) (
	*domain.Enrollment,
	error,
) {
	var enrollment domain.Enrollment

	query := "SELECT id, user_id, plan_id, effective_date, end_date FROM enrollments WHERE user_id = $1 AND now() BETWEEN effective_date AND end_date LIMIT 1"
	err := r.pg.Db.QueryRowContext(ctx, query, userID).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.PlanID,
		&enrollment.EffectiveDate,
		&enrollment.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoActivePlan
		}

		return nil, err
	}

	return &enrollment, nil
}

func (r *EnrollmentRepositoryImpl) FindActiveEnrollmentByUserIDTx(
	ctx context.Context,
	tx *sql.Tx,
	userID string,
	now time.Time,
) (
	*domain.Enrollment,
	error,
) {
	var enrollment domain.Enrollment

	query := "SELECT id, user_id, plan_id, effective_date, end_date FROM enrollments WHERE user_id = $1 AND $2 BETWEEN effective_date AND end_date LIMIT 1"
	err := tx.QueryRowContext(ctx, query, userID, now).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.PlanID,
		&enrollment.EffectiveDate,
		&enrollment.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoActivePlan
		}

		return nil, err
	}

	return &enrollment, nil
}
