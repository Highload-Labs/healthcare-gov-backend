package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/redis/go-redis/v9"
)

type PlanRepository interface {
	CountByState(ctx context.Context, state string) (int64, error)
	FindByState(ctx context.Context, state string, limit int, offset int) ([]domain.Plan, error)
	FindById(ctx context.Context, id string) (*domain.Plan, error)
}

var ErrPlanNotFound = errors.New("plan not found")

var planCacheKey = "plan:id:"

type PlanRepositoryImpl struct {
	postgres  *infra.Postgresql
	redisConn *redis.Client
}

func NewPlanRepository(postgres *infra.Postgresql, redisConn *redis.Client) PlanRepository {
	return &PlanRepositoryImpl{postgres: postgres, redisConn: redisConn}
}

func (r *PlanRepositoryImpl) CountByState(ctx context.Context, state string) (int64, error) {
	var count int64

	err := r.postgres.Db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) as total_data FROM plans WHERE state = $1",
		state,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PlanRepositoryImpl) FindByState(ctx context.Context, state string, limit int, offset int) (
	[]domain.Plan,
	error,
) {
	var plans []domain.Plan

	rows, err := r.postgres.Db.QueryContext(
		ctx,
		"SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE state = $1 ORDER BY id ASC LIMIT $2 OFFSET $3",
		state,
		limit,
		offset,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	for rows.Next() {
		var plan domain.Plan
		err = rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Provider,
			&plan.Tier,
			&plan.MonthlyPremium,
			&plan.Deductible,
			&plan.OutOfPocket,
			&plan.State,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func (r *PlanRepositoryImpl) FindById(ctx context.Context, id string) (*domain.Plan, error) {
	cacheKey := planCacheKey + id
	var plan domain.Plan

	val, err := r.redisConn.Get(ctx, cacheKey).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(val), &plan); err == nil {
			return &plan, err
		}
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	err = r.postgres.Db.QueryRowContext(
		ctx,
		"SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE id = $1",
		id,
	).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Provider,
		&plan.Tier,
		&plan.MonthlyPremium,
		&plan.Deductible,
		&plan.OutOfPocket,
		&plan.State,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	//nolint:gosec // cache population intentionally detached from request lifecycle
	go func(p domain.Plan) {
		backgroundCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		jsonData, err := json.Marshal(p)
		if err != nil {
			return
		}

		_ = r.redisConn.Set(backgroundCtx, "plan:id:"+p.ID, jsonData, 1*time.Hour).Err()
	}(plan)

	return &plan, nil
}
