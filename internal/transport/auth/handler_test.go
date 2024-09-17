package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth(t *testing.T) {
	storage := mocks.NewAuthStorage(t)
	authHandler := NewAuthHandler(storage, "encription-key")

	app := fiber.New()
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	t.Run("test Register", func(t *testing.T) {
		storage.On("Exists", mock.Anything, mock.Anything).Return(false, nil).Once()
		storage.On("Add", mock.Anything, mock.Anything).Return("1", nil).Once()

		body := strings.NewReader(`{
			"email": "any string",
			"name": "any string",
			"password": "any string"
		}`)

		req := httptest.NewRequest("POST", "/register", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("test Register again", func(t *testing.T) {
		storage.On("Exists", mock.Anything, mock.Anything).Return(true, nil).Once()

		body := strings.NewReader(`{
			"email": "any string",
			"name": "any string",
			"password": "any string"
		}`)

		req := httptest.NewRequest("POST", "/register", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test Register bad json", func(t *testing.T) {
		body := strings.NewReader(`{`)

		req := httptest.NewRequest("POST", "/register", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test Register empty json", func(t *testing.T) {
		body := strings.NewReader(`{}`)

		req := httptest.NewRequest("POST", "/register", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test Login successfully", func(t *testing.T) {
		user := &models.User{
			Email:    "email",
			Login:    "name",
			Password: "password",
		}

		storage.On("Get", mock.Anything, "login").Return(user, nil).Once()

		body := strings.NewReader(`{
			"email": "email",
			"login": "login",
			"password": "password"
		}`)

		req := httptest.NewRequest("POST", "/login", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("test Login with wrong password", func(t *testing.T) {
		user := &models.User{
			Email:    "email",
			Login:    "name",
			Password: "password",
		}

		storage.On("Get", mock.Anything, "login").Return(user, nil).Once()

		body := strings.NewReader(`{
			"email": "email",
			"login": "login",
			"password": "wrong password"
		}`)

		req := httptest.NewRequest("POST", "/login", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test Login with unknown login", func(t *testing.T) {
		storage.On("Get", mock.Anything, "unknown login").Return(nil, fmt.Errorf("%w", models.UserNotFound)).Once()

		body := strings.NewReader(`{
			"email": "email",
			"login": "unknown login",
			"password": "password"
		}`)

		req := httptest.NewRequest("POST", "/login", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("test Login bad json", func(t *testing.T) {
		body := strings.NewReader(`{`)

		req := httptest.NewRequest("POST", "/login", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test Login empty json", func(t *testing.T) {
		body := strings.NewReader(`{}`)

		req := httptest.NewRequest("POST", "/login", body)
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
