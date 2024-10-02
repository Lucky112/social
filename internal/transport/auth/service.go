package auth

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

// Сервис зарегистрированных пользователей
type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	NewUser(ctx context.Context, user *models.User) (string, error)
}
