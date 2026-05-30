package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

type PlanRepository interface {
	FindByState(ctx context.Context, state string) ([]domain.Plan, error)
	FindById(ctx context.Context, id string) (*domain.Plan, error)
}

var ErrPlanNotFound = errors.New("plan not found")

type PlanRepositoryImpl struct {
	postgres *infra.Postgresql
}

func NewPlanRepository(postgres *infra.Postgresql) PlanRepository {
	return &PlanRepositoryImpl{postgres: postgres}
}

func (r *PlanRepositoryImpl) FindByState(ctx context.Context, state string) ([]domain.Plan, error) {
	var plans []domain.Plan

	rows, err := r.postgres.Db.QueryContext(ctx, "SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE state = $1", state)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	for rows.Next() {
		var plan domain.Plan
		err = rows.Scan(&plan.ID, &plan.Name, &plan.Provider, &plan.Tier, &plan.MonthlyPremium, &plan.Deductible, &plan.OutOfPocket, &plan.State, &plan.CreatedAt, &plan.UpdatedAt)
		if err != nil {
			return nil, err
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func (r *PlanRepositoryImpl) FindById(ctx context.Context, id string) (*domain.Plan, error) {
	var plan domain.Plan

	err := r.postgres.Db.QueryRowContext(ctx, "SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE id = $1", id).Scan(&plan.ID, &plan.Name, &plan.Provider, &plan.Tier, &plan.MonthlyPremium, &plan.Deductible, &plan.OutOfPocket, &plan.State, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	return &plan, nil
}
