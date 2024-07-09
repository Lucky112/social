package auth

import "github.com/Lucky112/social/internal/models"

// Хранилище зарегистрированных пользователей
type AuthStorage interface {
	Exists(userId string) bool
	Get(userId string) (*models.User, error)
	Add(userId string, user *models.User) error
}
