package auth

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

// Хранилище зарегистрированных пользователей
type AuthStorage interface {
	Exists(ctx context.Context, user *models.User) (bool, error)
	Get(ctx context.Context, userId string) (*models.User, error)
	Add(ctx context.Context, user *models.User) (string, error)
}
