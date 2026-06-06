package infra

import (
	"errors"
	"log/slog"
	"os"

	healthcare_gov_backend "github.com/Highload-Labs/healthcare-gov-backend"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunDatabaseMigrations(pg *Postgresql) {
	sourceDriver, err := iofs.New(healthcare_gov_backend.MigrationFiles, "migrations")
	if err != nil {
		slog.Error("Failed to initialize migration source driver", "error", err)
		os.Exit(1)
	}

	dbInstance := pg.Db
	migrationTargetDriver, err := postgres.WithInstance(dbInstance, &postgres.Config{})
	if err != nil {
		slog.Error("Failed to initialize migration database driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", migrationTargetDriver)
	if err != nil {
		slog.Error("Failed to create migrate instance", "error", err)
		os.Exit(1)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database schema is already up to date. No changes applied.")
		} else {
			slog.Error("Migration execution failed", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Info("Database migrations applied successfully!")
	}
}
