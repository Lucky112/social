package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth(t *testing.T) {
	storage := mocks.NewAuthStorage(t)
	authService := NewAuthService(storage)

	t.Run("test NewUser", func(t *testing.T) {
		storage.On("Exists", mock.Anything, mock.Anything).Return(false, nil).Once()
		storage.On("Add", mock.Anything, mock.Anything).Return("1", nil).Once()

		user := &models.User{
			Email:    "email",
			Login:    "login",
			Password: "pwd",
		}

		id, err := authService.NewUser(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, "1", id)
	})

	t.Run("test Register again", func(t *testing.T) {
		storage.On("Exists", mock.Anything, mock.Anything).Return(true, nil).Once()

		user := &models.User{
			Email:    "email",
			Login:    "login",
			Password: "pwd",
		}

		_, err := authService.NewUser(context.Background(), user)
		assert.ErrorIs(t, err, models.UserAlreadyExists)
	})

	t.Run("test Login successfully", func(t *testing.T) {
		userId := "2"
		login := "login"
		password := "pwd"
		hashedPwd, _ := hashAndSalt([]byte(password))

		user := &models.User{
			Id:             userId,
			Email:          "email",
			Login:          login,
			HashedPassword: hashedPwd,
		}

		storage.On("Get", mock.Anything, "login").Return(user, nil).Once()

		id, err := authService.Login(context.Background(), login, password)
		assert.NoError(t, err)
		assert.Equal(t, userId, id)
	})

	t.Run("test Login with wrong password", func(t *testing.T) {
		userId := "2"
		login := "login"
		password := "pwd"
		hashedPwd, _ := hashAndSalt([]byte(password))

		user := &models.User{
			Id:             userId,
			Email:          "email",
			Login:          login,
			HashedPassword: hashedPwd,
		}

		storage.On("Get", mock.Anything, "login").Return(user, nil).Once()

		_, err := authService.Login(context.Background(), login, "wrong password")
		assert.ErrorIs(t, err, models.UserBadCredentials)
	})

	t.Run("test Login with unknown login", func(t *testing.T) {
		login := "unknown login"

		storage.On("Get", mock.Anything, login).Return(nil, fmt.Errorf("%w", models.UserNotFound)).Once()

		_, err := authService.Login(context.Background(), login, "password")
		assert.ErrorIs(t, err, models.UserNotFound)
	})

	t.Run("test NewUser storage errors", func(t *testing.T) {
		user := &models.User{
			Email:    "email",
			Login:    "login",
			Password: "pwd",
		}

		storage.On("Exists", mock.Anything, mock.Anything).Return(false, errors.New("storage error")).Once()

		_, err := authService.NewUser(context.Background(), user)
		assert.Error(t, err)

		storage.On("Exists", mock.Anything, mock.Anything).Return(false, nil).Once()
		storage.On("Add", mock.Anything, mock.Anything).Return("1", errors.New("storage error")).Once()

		_, err = authService.NewUser(context.Background(), user)
		assert.Error(t, err)
	})

	t.Run("test Login storage error", func(t *testing.T) {
		login := "login"

		storage.On("Get", mock.Anything, login).Return(nil, errors.New("storage error")).Once()

		_, err := authService.Login(context.Background(), login, "password")
		assert.Error(t, err)
	})
}
