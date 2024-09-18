package auth

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/transport/jwt"
)

// Обработчик HTTP-запросов на регистрацию и аутентификацию пользователей
type AuthHandler struct {
	service  AuthService
	jwtKey   []byte
	validate *validator.Validate
}

func NewAuthHandler(service AuthService, jwtKey []byte) AuthHandler {
	return AuthHandler{
		service:  service,
		jwtKey:   jwtKey,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Обработчик HTTP-запросов на регистрацию пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	regReq := registerRequest{}
	err := c.BodyParser(&regReq)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			registerError{fmt.Sprintf("failed to parse body: %v", err)},
		)
		return nil
	}

	err = h.validate.Struct(regReq)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			registerError{fmt.Sprintf("invalid body: %v", err)},
		)
		return nil
	}

	user := &models.User{
		Email:    regReq.Email,
		Login:    regReq.Login,
		Password: regReq.Password,
	}

	id, err := h.service.NewUser(c.Context(), user)
	if err != nil {
		if errors.Is(err, models.UserAlreadyExists) {
			c.Status(fiber.StatusBadRequest).JSON(
				registerError{"the user for given email or login already exists"},
			)
		} else {
			c.Status(fiber.StatusInternalServerError).JSON(
				registerError{fmt.Sprintf("failed to create new user: %v", err)},
			)
		}

		return nil
	}

	err = c.Status(fiber.StatusCreated).JSON(
		registerResponse{id},
	)
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}

// Обработчик HTTP-запросов на вход в аккаунт
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	loginReq := loginRequest{}
	err := c.BodyParser(&loginReq)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			loginError{fmt.Sprintf("failed to parse body: %v", err)},
		)
		return nil
	}

	err = h.validate.Struct(loginReq)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			loginError{fmt.Sprintf("invalid body: %v", err)},
		)
		return nil
	}

	userId, err := h.service.Login(c.Context(), loginReq.Login, loginReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.UserNotFound):
			c.Status(fiber.StatusNotFound).JSON(
				loginError{"the user for given login not found"},
			)
		case errors.Is(err, models.UserBadCredentials):
			c.Status(fiber.StatusBadRequest).JSON(
				loginError{"login or password is incorrect"},
			)
		default:
			c.Status(fiber.StatusInternalServerError).JSON(
				loginError{fmt.Sprintf("failed to login: %v", err)},
			)
		}

		return nil
	}

	token, err := jwt.MakeToken(userId, h.jwtKey)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(
			loginError{fmt.Sprintf("failed to create JWT-token: %v", err)},
		)
		return nil
	}

	err = c.JSON(loginResponse{AccessToken: token})
	if err != nil {
		return fmt.Errorf("sending response: %v", err)
	}

	return nil
}
