package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage UsersStorage
}

// Хранилище зарегистрированных пользователей
type UsersStorage interface {
	Exists(ctx context.Context, user *models.User) (bool, error)
	Get(ctx context.Context, userId string) (*models.User, error)
	Add(ctx context.Context, user *models.User) (string, error)
}

func NewAuthService(storage UsersStorage) AuthService {
	return AuthService{
		storage: storage,
	}
}

func (s AuthService) NewUser(ctx context.Context, user *models.User) (string, error) {
	hashedPassword, err := hashAndSalt([]byte(user.Password))
	if err != nil {
		return "", fmt.Errorf("hashing password: %v", err)
	}
	user.HashedPassword = hashedPassword

	exists, err := s.storage.Exists(ctx, user)
	if err != nil {
		return "", fmt.Errorf("checking user existence: %v", err)
	}
	if exists {
		return "", models.UserAlreadyExists
	}

	id, err := s.storage.Add(ctx, user)
	if err != nil {
		return "", fmt.Errorf("creating new user: %v", err)
	}

	return id, nil
}

func (s AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.Get(ctx, login)
	if err != nil {
		if errors.Is(err, models.UserNotFound) {
			return "", err
		}

		return "", fmt.Errorf("looking for user '%s': %v", login, err)
	}

	err = checkHash([]byte(password), user.HashedPassword)
	if err != nil {
		return "", models.UserBadCredentials
	}

	return user.Id, nil
}

func hashAndSalt(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %v", err)
	}

	return hash, nil
}

func checkHash(pwd, hashedPwd []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPwd, pwd)
	if err != nil {
		return fmt.Errorf("comparing passwords: %v", err)
	}

	return nil
}
