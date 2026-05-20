package infra

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
)

type Postgresql struct {
	Db *sql.DB
}

type pgConnParam struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DBName  string
	SSLMode string
}

func NewPostgresql(cfg *config.Config) *Postgresql {
	return &Postgresql{
		Db: CreatePostgresqlConnection(pgConnParam{
			Host:    cfg.DatabaseHost,
			Port:    cfg.DatabasePort,
			User:    cfg.DatabaseUser,
			Pass:    cfg.DatabasePass,
			DBName:  cfg.DatabaseName,
			SSLMode: cfg.DatabaseSSLMode,
		}),
	}
}

func CreatePostgresqlConnection(param pgConnParam) *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", param.User, param.Pass, param.DBName, param.SSLMode)

	if param.Pass == "" {
		connStr = fmt.Sprintf("user=%s dbname=%s sslmode=%s", param.User, param.DBName, param.SSLMode)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("database error", "error", err.Error())
		return nil
	}

	return db
}
