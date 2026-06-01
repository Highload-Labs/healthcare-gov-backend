package infra

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	ctx := context.Background()

	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	rdb := redis.NewClient(
		&redis.Options{
			Addr: addr,
			DB:   cfg.RedisDB,
		},
	)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.Error("redis error", "redis ping failed", err.Error())
	}

	slog.Info("redis", "redis connect success")

	return rdb
}
