package profiles

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/transport/jwt"
	"github.com/Lucky112/social/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProfiles(t *testing.T) {
	service := mocks.NewProfilesService(t)
	profilesHandler := NewProfilesHandler(service)
	signingKey := []byte("signing-key")

	app := fiber.New()
	app.Use(jwt.Middleware(signingKey))
	app.Post("/profiles", profilesHandler.CreateProfile)
	app.Get("/profiles", profilesHandler.GetProfiles)
	app.Get("/profiles/:id", profilesHandler.GetProfile)

	t.Run("test CreateProfile", func(t *testing.T) {
		userId := "1"
		profileId := "23"

		service.On("Add", mock.Anything, mock.Anything).Return(profileId, nil).Once()

		body := strings.NewReader(`{
			"name": "Alfred",
			"surname": "Winner",
			"sex": "male",
			"age": 34,
			"hobbies": "cycling, reading, chess, parties",
			"address": "Moscow"
		}`)

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/profiles", body)
		req.Header.Add("Content-type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("test CreateProfile failed", func(t *testing.T) {
		userId := "1"

		service.On("Add", mock.Anything, mock.Anything).Return("", fmt.Errorf("service error")).Once()

		body := strings.NewReader(`{
			"name": "Alfred",
			"surname": "Winner",
			"sex": "male",
			"age": 34,
			"hobbies": "cycling, reading, chess, parties",
			"address": "Moscow"
		}`)

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/profiles", body)
		req.Header.Add("Content-type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("test CreateProfile bad profile", func(t *testing.T) {
		userId := "1"

		body := strings.NewReader(`{
			"name": "Alfred",
			"surname": "Winner",
			"sex": "malformed sex",
			"age": 34,
			"hobbies": "cycling, reading, chess, parties",
			"address": "Moscow"
		}`)

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/profiles", body)
		req.Header.Add("Content-type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test CreateProfile bad json", func(t *testing.T) {
		userId := "1"
		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		body := strings.NewReader(`{`)

		req := httptest.NewRequest("POST", "/profiles", body)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test CreateProfile empty json", func(t *testing.T) {
		userId := "1"
		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		body := strings.NewReader(`{}`)

		req := httptest.NewRequest("POST", "/profiles", body)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("test GetProfile", func(t *testing.T) {
		userId := "1"
		profileId := "23"

		profile := &models.Profile{
			Name: "username",
		}

		service.On("Get", mock.Anything, profileId).Return(profile, nil).Once()

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", fmt.Sprintf("/profiles/%s", profileId), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("test GetProfile not found", func(t *testing.T) {
		userId := "1"
		profileId := "23"

		service.On("Get", mock.Anything, profileId).Return(nil, fmt.Errorf("%w", models.ProfileNotFound)).Once()

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", fmt.Sprintf("/profiles/%s", profileId), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("test GetProfile failed", func(t *testing.T) {
		userId := "1"
		profileId := "23"

		service.On("Get", mock.Anything, profileId).Return(nil, fmt.Errorf("service error")).Once()

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", fmt.Sprintf("/profiles/%s", profileId), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("test GetProfiles", func(t *testing.T) {
		userId := "1"

		profiles := []*models.Profile{
			{
				Name: "username",
			},
		}

		service.On("GetAll", mock.Anything).Return(profiles, nil).Once()

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/profiles", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("test GetProfiles failed", func(t *testing.T) {
		userId := "1"

		service.On("GetAll", mock.Anything).Return(nil, fmt.Errorf("service error")).Once()

		token, err := jwt.MakeToken(userId, signingKey)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/profiles", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
