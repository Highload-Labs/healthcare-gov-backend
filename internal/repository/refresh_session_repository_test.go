package repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestRefreshSessionRepository_Create(t *testing.T) {
	redisConn, mock := redismock.NewClientMock()

	fixedExpiry := time.Now().Add(time.Hour)
	expectedTTL := time.Until(fixedExpiry)

	userID := "1"
	tokenHash := "test-token-hash"
	expectedKey := refreshSessionCacheKey + tokenHash

	mock.ExpectSet(expectedKey, userID, expectedTTL).SetVal("OK")

	repo := &RefreshSessionRepositoryImpl{
		redisConn: redisConn,
	}

	err := repo.Create(context.Background(), userID, tokenHash, fixedExpiry)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating a refresh session", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
