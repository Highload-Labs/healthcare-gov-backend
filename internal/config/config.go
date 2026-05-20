package config

import (
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	ServerPort string
	BcryptCost int

	DatabaseHost    string
	DatabasePort    string
	DatabaseUser    string
	DatabasePass    string
	DatabaseName    string
	DatabaseSSLMode string

	AccessTokenExpired   time.Duration
	RefreshTokenExpired  time.Duration
	JwtAccessSigningKey  []byte
	JwtRefreshSigningKey []byte
}

var config *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(
		func() {
			err := godotenv.Load()
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}

			serverPort := os.Getenv("SERVER_PORT")
			if serverPort == "" {
				serverPort = "8080"
			}

			bcryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
			if err != nil {
				slog.Error("configuration error", "details", err.Error())
				os.Exit(1)
			}

			if bcryptCost == 0 {
				bcryptCost = bcrypt.DefaultCost
			}

			databaseHost := os.Getenv("DATABASE_HOST")
			if databaseHost == "" {
				slog.Error("configuration error", "details", "unable to determine database host")
				os.Exit(1)
			}

			databasePost := os.Getenv("DATABASE_PORT")
			if databasePost == "" {
				slog.Error("configuration error", "details", "unable to determine database port")
				os.Exit(1)
			}

			databaseUser := os.Getenv("DATABASE_USER")
			if databaseUser == "" {
				slog.Error("configuration error", "details", "unable to determine database user")
				os.Exit(1)
			}

			databasePass := os.Getenv("DATABASE_PASS")

			databaseName := os.Getenv("DATABASE_NAME")
			if databaseName == "" {
				slog.Error("configuration error", "details", "unable to determine database name")
				os.Exit(1)
			}

			databaseSSLMode := os.Getenv("DATABASE_SSLMODE")
			if databaseSSLMode == "" {
				databaseSSLMode = "disable"
			}

			accessTokenExpired, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED"))
			if err != nil {
				slog.Error("configuration error", "details", err.Error())
				os.Exit(1)
			}

			if accessTokenExpired == 0 {
				accessTokenExpired = 1 * time.Hour
			}

			refreshTokenExpired, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED"))
			if err != nil {
				slog.Error("configuration error", "details", err.Error())
				os.Exit(1)
			}

			if refreshTokenExpired == 0 {
				refreshTokenExpired = 24 * time.Hour
			}

			jwtAccessSigningKey := os.Getenv("JWT_ACCESS_SIGNING_KEY")
			if jwtAccessSigningKey == "" {
				slog.Error("configuration error", "details", "empty jwt access signing key")
				os.Exit(1)
			}

			jwtRefreshSigningKey := os.Getenv("JWT_REFRESH_SIGNING_KEY")
			if jwtAccessSigningKey == "" {
				slog.Error("configuration error", "details", "empty jwt refresh signing key")
				os.Exit(1)
			}

			config = &Config{
				ServerPort:           serverPort,
				BcryptCost:           bcryptCost,
				DatabaseHost:         databaseHost,
				DatabasePort:         databasePost,
				DatabaseUser:         databaseUser,
				DatabasePass:         databasePass,
				DatabaseName:         databaseName,
				DatabaseSSLMode:      databaseSSLMode,
				AccessTokenExpired:   accessTokenExpired,
				RefreshTokenExpired:  refreshTokenExpired,
				JwtAccessSigningKey:  []byte(jwtAccessSigningKey),
				JwtRefreshSigningKey: []byte(jwtRefreshSigningKey),
			}
		},
	)

	return config
}
