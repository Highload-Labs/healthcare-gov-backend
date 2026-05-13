package memory

import (
	"context"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"sync"
)

type userRepository struct {
	mu    sync.RWMutex
	users map[string]domain.User // key is email
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{
		users: make(map[string]domain.User),
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return &user, nil
}
