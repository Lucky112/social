package storage

import (
	"fmt"

	"github.com/Lucky112/social/internal/models"
)

type AuthStorage map[string]*models.User

func NewAuthStorage() AuthStorage {
	return make(AuthStorage)
}

func (a AuthStorage) Exists(userId string) bool {
	_, exists := a[userId]
	return exists
}

func (a AuthStorage) Add(userId string, user *models.User) error {
	a[userId] = user
	return nil
}

func (a AuthStorage) Get(userID string) (*models.User, error) {
	user, exists := a[userID]
	if !exists {
		return nil, fmt.Errorf("looking '%s' up: %w", userID, models.UserNotFound)
	}

	return user, nil
}
