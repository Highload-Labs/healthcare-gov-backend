package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

func TestEnrollmentRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta("INSERT INTO enrollments (user_id, plan_id, effective_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id")

	now := time.Now()
	rows := mock.NewRows([]string{"id"}).AddRow("1")

	pg := &infra.Postgresql{Db: db}
	repo := &EnrollmentRepositoryImpl{pg}

	mock.ExpectQuery(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

	id, err := repo.Create(context.Background(), domain.Enrollment{
		UserID:        "1",
		PlanID:        "1",
		EffectiveDate: now,
		EndDate:       now,
	})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating an enrollment", err)
	}

	if id != "1" {
		t.Errorf("expected id '1', got '%s'", id)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
