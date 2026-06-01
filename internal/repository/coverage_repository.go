package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/redis/go-redis/v9"
)

var ErrCoverageNotFound = errors.New("coverage not found")

var coverageCacheKey = "coverage:zip:"

type CoverageRepository interface {
	FindByZipcode(ctx context.Context, zipcode string) (*domain.Coverage, error)
}

type CoverageRepositoryImpl struct {
	postgres  *infra.Postgresql
	redisConn *redis.Client
}

func NewCoverageRepository(postgres *infra.Postgresql, redisConn *redis.Client) CoverageRepository {
	return &CoverageRepositoryImpl{postgres: postgres, redisConn: redisConn}
}

func (r *CoverageRepositoryImpl) FindByZipcode(ctx context.Context, zipcode string) (*domain.Coverage, error) {
	cacheKey := coverageCacheKey + zipcode

	stateVal, err := r.redisConn.Get(ctx, cacheKey).Result()
	if err == nil {
		return &domain.Coverage{State: stateVal}, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var coverage domain.Coverage

	query := "SELECT id, state, zipcode_start, zipcode_end FROM coverages WHERE $1 BETWEEN zipcode_start AND zipcode_end LIMIT 1"
	err = r.postgres.Db.QueryRowContext(ctx, query, zipcode).Scan(
		&coverage.ID,
		&coverage.State,
		&coverage.ZipcodeStart,
		&coverage.ZipcodeEnd,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCoverageNotFound
		}

		return nil, err
	}

	go func(zip string, state string) {
		backgroundCtx := context.Background()
		err = r.redisConn.Set(backgroundCtx, cacheKey, state, 0).Err()
		if err != nil {
			return
		}
	}(zipcode, coverage.State)

	return &coverage, nil
}
