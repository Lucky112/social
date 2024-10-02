package inmemory

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/internal/models"
)

type AuthStorage map[string]*models.User

func NewAuthStorage() AuthStorage {
	return make(AuthStorage)
}

func (a AuthStorage) Exists(ctx context.Context, user *models.User) (bool, error) {
	_, exists := a[user.Id]
	return exists, nil
}

func (a AuthStorage) Add(ctx context.Context, user *models.User) (string, error) {
	userId := generateId()
	a[userId] = user

	return userId, nil
}

func (a AuthStorage) Get(ctx context.Context, userID string) (*models.User, error) {
	user, exists := a[userID]
	if !exists {
		return nil, fmt.Errorf("looking '%s' up: %w", userID, models.UserNotFound)
	}

	return user, nil
}
