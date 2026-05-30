package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

func TestPlanRepository_FindByState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta("SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE state = $1")

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "provider", "tier", "monthly_premium", "deductible", "out_of_pocket_max", "state", "created_at", "updated_at"}).AddRow("1", "test", "test", "bronze", 50.00, 50.00, 50.00, "test", now, now)

	mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)

	pg := &infra.Postgresql{
		Db: db,
	}

	repo := &PlanRepositoryImpl{
		postgres: pg,
	}

	plans, err := repo.FindByState(context.Background(), "test")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding plans by state", err)
	}

	if len(plans) != 1 {
		t.Fatalf("there should be 1 plan")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
