package config

import (
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	ServerPort string
	BcryptCost int
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
				slog.Error(err.Error())
				os.Exit(1)
			}

			if bcryptCost == 0 {
				bcryptCost = bcrypt.DefaultCost
			}

			config = &Config{
				ServerPort: serverPort,
				BcryptCost: bcryptCost,
			}
		},
	)

	return config
}
