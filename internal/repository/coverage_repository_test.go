package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/go-redis/redismock/v9"
)

func TestCoverageRepository_FindByZipcode_CacheHit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	redisConn, redisMock := redismock.NewClientMock()

	pg := &infra.Postgresql{
		Db: db,
	}

	zipcode := "1"
	cacheKey := coverageCacheKey + zipcode

	redisMock.ExpectGet(cacheKey).SetVal("California")

	repo := NewCoverageRepository(pg, redisConn)

	//rows := sqlmock.NewRows([]string{"id", "state", "zipcode_start", "zipcode_end"}).AddRow("test", "test", "1", "1")
	//
	//query := regexp.QuoteMeta(`SELECT id, state, zipcode_start, zipcode_end FROM coverages WHERE $1 BETWEEN zipcode_start AND zipcode_end LIMIT 1`)
	//mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)

	coverage, err := repo.FindByZipcode(context.Background(), "1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding coverage by zipcode", err)
	}

	if coverage.State != "California" {
		t.Fatalf("wanted California, got %s", coverage.State)
	}

	if err = redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
