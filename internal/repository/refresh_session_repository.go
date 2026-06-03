package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session not found")

var ErrInvalidTTL = errors.New("invalid TTL")

var refreshSessionCacheKey = "session:refresh:"

type RefreshSessionRepository interface {
	Create(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error
	Revoke(ctx context.Context, tokenHash string) error
	FindByHash(ctx context.Context, tokenHash string) (*domain.RefreshSession, error)
}

type RefreshSessionRepositoryImpl struct {
	redisConn *redis.Client
}

func NewRefreshTokenRepository(redisConn *redis.Client) RefreshSessionRepository {
	return &RefreshSessionRepositoryImpl{redisConn: redisConn}
}

func (r *RefreshSessionRepositoryImpl) Create(
	ctx context.Context,
	userID, tokenHash string,
	expiresAt time.Time,
) error {
	ttl := time.Until(expiresAt)

	if ttl < 0 {
		return ErrInvalidTTL
	}

	cacheKey := refreshSessionCacheKey + tokenHash

	_, err := r.redisConn.Set(ctx, cacheKey, userID, ttl).Result()

	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshSessionRepositoryImpl) Revoke(
	ctx context.Context,
	tokenHash string,
) error {
	_, err := r.redisConn.Del(ctx, refreshSessionCacheKey+tokenHash).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshSessionRepositoryImpl) FindByHash(ctx context.Context, tokenHash string) (
	*domain.RefreshSession,
	error,
) {
	session := domain.RefreshSession{}

	cacheKey := refreshSessionCacheKey + tokenHash

	hashVal, err := r.redisConn.Get(ctx, cacheKey).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotFound
	}

	session.TokenHash = tokenHash
	session.UserID = hashVal

	return &session, nil
}
