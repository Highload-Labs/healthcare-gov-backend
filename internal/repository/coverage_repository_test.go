package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

func TestCoverageRepository_FindByZipcode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pg := &infra.Postgresql{
		Db: db,
	}

	repo := NewCoverageRepository(pg)

	rows := sqlmock.NewRows([]string{"id", "state", "zipcode_start", "zipcode_end"}).AddRow("test", "test", "1", "1")

	query := regexp.QuoteMeta(`SELECT id, state, zipcode_start, zipcode_end FROM coverages WHERE $1 BETWEEN zipcode_start AND zipcode_end LIMIT 1`)
	mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)

	coverage, err := repo.FindByZipcode(context.Background(), "1")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding coverage by zipcode", err)
	}

	if coverage.ID != "test" {
		t.Fatalf("expected id 'test', got '%s'", coverage.ID)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
