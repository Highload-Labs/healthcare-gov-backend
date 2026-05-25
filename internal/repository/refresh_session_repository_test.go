package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

func TestRefreshSessionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pg := &infra.Postgresql{
		Db: db,
	}

	fixedExpiry := time.Now().Add(time.Hour)

	query := regexp.QuoteMeta("INSERT INTO refresh_sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)")
	mock.ExpectExec(query).WithArgs(
		"1",
		"test-token-hash",
		fixedExpiry,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &RefreshSessionRepositoryImpl{
		postgres: pg,
	}

	err = repo.Create(context.Background(), "1", "test-token-hash", fixedExpiry)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating a refresh session", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
