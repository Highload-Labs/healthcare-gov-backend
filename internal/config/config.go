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
	ServerPort           string
	BcryptCost           int
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
				AccessTokenExpired:   accessTokenExpired,
				RefreshTokenExpired:  refreshTokenExpired,
				JwtAccessSigningKey:  []byte(jwtAccessSigningKey),
				JwtRefreshSigningKey: []byte(jwtRefreshSigningKey),
			}
		},
	)

	return config
}
