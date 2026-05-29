package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

var ErrCoverageNotFound = errors.New("coverage not found")

type CoverageRepository interface {
	FindByZipcode(ctx context.Context, zipcode string) (*domain.Coverage, error)
}

type CoverageRepositoryImpl struct {
	postgres *infra.Postgresql
}

func NewCoverageRepository(postgres *infra.Postgresql) CoverageRepository {
	return &CoverageRepositoryImpl{postgres: postgres}
}

func (r *CoverageRepositoryImpl) FindByZipcode(ctx context.Context, zipcode string) (*domain.Coverage, error) {
	var coverage domain.Coverage

	query := "SELECT id, state, zipcode_start, zipcode_end FROM coverages WHERE $1 BETWEEN zipcode_start AND zipcode_end LIMIT 1"
	err := r.postgres.Db.QueryRowContext(ctx, query, zipcode).Scan(&coverage.Id, &coverage.State, &coverage.ZipcodeStart, &coverage.ZipcodeEnd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCoverageNotFound
		}

		return nil, err
	}

	return &coverage, nil
}
