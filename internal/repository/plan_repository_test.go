package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/go-redis/redismock/v9"
	"github.com/prometheus/client_golang/prometheus"
)

func TestPlanRepository_CountByState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta("SELECT COUNT(*) as total_data FROM plans WHERE state = $1")

	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	pg := &infra.Postgresql{
		Db: db,
	}

	repo := &PlanRepositoryImpl{
		postgres: pg,
	}

	count, err := repo.CountByState(context.Background(), "state")
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("want 1, got %d", count)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlanRepository_FindByState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta("SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE state = $1 ORDER BY id ASC LIMIT $2 OFFSET $3")

	now := time.Now()

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"name",
			"provider",
			"tier",
			"monthly_premium",
			"deductible",
			"out_of_pocket_max",
			"state",
			"created_at",
			"updated_at",
		},
	).AddRow("1", "test", "test", "bronze", 50.00, 50.00, 50.00, "test", now, now)

	mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

	pg := &infra.Postgresql{
		Db: db,
	}

	repo := &PlanRepositoryImpl{
		postgres: pg,
	}

	plans, err := repo.FindByState(context.Background(), "test", 1, 0)
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

func TestPlanRepository_FindById_CacheHit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	redisConn, redisMock := redismock.NewClientMock()

	//query := regexp.QuoteMeta("SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE id = $1")
	//
	//now := time.Now()
	//rows := sqlmock.NewRows(
	//	[]string{
	//		"id",
	//		"name",
	//		"provider",
	//		"tier",
	//		"monthly_premium",
	//		"deductible",
	//		"out_of_pocket_max",
	//		"state",
	//		"created_at",
	//		"updated_at",
	//	},
	//).AddRow("1", "test", "test", "bronze", 50.00, 50.00, 50.00, "test", now, now)
	//
	//mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)

	pg := &infra.Postgresql{
		Db: db,
	}

	metrics := infra.NewMetrics(prometheus.DefaultRegisterer)

	repo := &PlanRepositoryImpl{
		postgres:  pg,
		redisConn: redisConn,
		metrics:   metrics,
	}

	id := "1"
	cacheKey := planCacheKey + id
	redisMock.ExpectGet(cacheKey).SetVal(`{"id":"1"}`)

	plan, err := repo.FindById(context.Background(), id)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding plans by id", err)
	}

	if plan.ID != id {
		t.Fatalf("id expected to be 1")
	}

	if err = redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
