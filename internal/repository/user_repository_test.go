package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pg := &infra.Postgresql{
		Db: db,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	query := regexp.QuoteMeta("INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id")
	mock.ExpectQuery(query).WithArgs(
		"test@gmail.com",
		"test",
		"test",
	).WillReturnRows(rows)

	repo := NewUserRepository(pg)
	ctx := context.Background()

	userID, err := repo.Create(
		ctx, domain.User{
			Email:    "test@gmail.com",
			Username: "test",
			Password: "test",
		},
	)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating a user", err)
	}

	if userID != "1" {
		t.Errorf("got user id %s, expected 1", userID)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pg := &infra.Postgresql{
		Db: db,
	}

	repo := NewUserRepository(pg)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password"}).AddRow(
		"1",
		"test@gmail.com",
		"tester",
		"hashed-password",
	)

	query := regexp.QuoteMeta("SELECT id, email, username, password FROM users WHERE email = $1")
	mock.ExpectQuery(query).WithArgs("test@gmail.com").WillReturnRows(rows)
	user, err := repo.FindByEmail(ctx, "test@gmail.com")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding user", err)
	}

	if user.Email != "test@gmail.com" {
		t.Errorf("expected email test@gmail.com but got %s", user.Email)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pg := &infra.Postgresql{
		Db: db,
	}

	repo := NewUserRepository(pg)

	mock.ExpectQuery(
		`SELECT id, email, username, password FROM users WHERE email = \$1`,
	).
		WithArgs("unknown@gmail.com").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByEmail(context.Background(), "unknown@gmail.com")

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func setupBenchmarkPostgres(b *testing.B) *infra.Postgresql {
	cfg := &config.Config{
		DatabaseHost:    "localhost",
		DatabasePort:    "5432",
		DatabaseUser:    "postgres",
		DatabasePass:    "althea",
		DatabaseName:    "healthcare_gov_benchmark",
		DatabaseSSLMode: "disable",
	}

	pg := infra.NewPostgresql(cfg)
	if pg.Db == nil {
		b.Fatal("Database not initialized")
	}

	return pg
}

func BenchmarkUserRepository_Create(b *testing.B) {
	pg := setupBenchmarkPostgres(b)

	_, err := pg.Db.Exec("TRUNCATE TABLE users")
	if err != nil {
		b.Fatal(err)
	}

	repo := NewUserRepository(pg)
	ctx := context.Background()

	users := make([]domain.User, b.N)

	for i := 0; i < b.N; i++ {
		users[i] = domain.User{
			Email:    fmt.Sprintf("test%d@gmail.com", i),
			Username: "test",
			Password: "test",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid, err := repo.Create(
			ctx, users[i],
		)

		if err != nil {
			b.Fatal(err)
		}

		if uid == "" {
			b.Fatal("uid was empty")
		}
	}

	b.Cleanup(
		func() {
			_ = pg.Db.Close()
		},
	)
}

func BenchmarkUserRepository_FindByEmail(b *testing.B) {
	pg := setupBenchmarkPostgres(b)
	defer pg.Db.Close()

	_, err := pg.Db.Exec("TRUNCATE TABLE users")
	if err != nil {
		b.Fatal(err)
	}

	_, err = pg.Db.Exec(
		`
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
	`,
		"test@gmail.com",
		"tester",
		"hashed",
	)

	if err != nil {
		b.Fatal(err)
	}

	repo := NewUserRepository(pg)
	ctx := context.Background()
	b.ResetTimer()

	for b.Loop() {
		_, err = repo.FindByEmail(ctx, "test@gmail.com")
		if err != nil {
			b.Fatal(err)
		}
	}

	b.Cleanup(
		func() {
			_ = pg.Db.Close()
		},
	)
}

func BenchmarkUserRepository_FindByEmail_NotFound(b *testing.B) {
	pg := setupBenchmarkPostgres(b)

	_, err := pg.Db.Exec("TRUNCATE TABLE users")
	if err != nil {
		b.Fatal(err)
	}

	_, err = pg.Db.Exec(
		`
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
	`,
		"test@gmail.com",
		"tester",
		"hashed",
	)

	if err != nil {
		b.Fatal(err)
	}

	repo := NewUserRepository(pg)
	ctx := context.Background()
	b.ResetTimer()

	for b.Loop() {
		_, err = repo.FindByEmail(ctx, "asdas@gmail.com")
		if err != nil {
			if !errors.Is(err, ErrUserNotFound) {
				b.Fatal(err)
			}
		}
	}

	b.Cleanup(
		func() {
			_ = pg.Db.Close()
		},
	)
}
