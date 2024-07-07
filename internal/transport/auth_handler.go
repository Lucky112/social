package transport

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Lucky112/social/internal/models"
)

type (
	// Обработчик HTTP-запросов на регистрацию и аутентификацию пользователей
	AuthHandler struct {
		storage AuthStorage
		signKey []byte
	}

	// Хранилище зарегистрированных пользователей
	AuthStorage interface {
		Exists(userId string) bool
		Get(userId string) (*models.User, error)
		Add(userId string, user *models.User) error
	}
)

// Структура HTTP-запроса на регистрацию пользователя
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Структура HTTP-запроса на вход в аккаунт
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Структура HTTP-ответа на вход в аккаунт
// В ответе содержится JWT-токен авторизованного пользователя
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

var (
	errBadCredentials = errors.New("email or password is incorrect")
)

func NewAuthHandler(storage AuthStorage, signKey string) AuthHandler {
	return AuthHandler{
		storage: storage,
		signKey: []byte(signKey),
	}
}

// Обработчик HTTP-запросов на регистрацию пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	regReq := RegisterRequest{}
	err := c.BodyParser(&regReq)
	if err != nil {
		return fmt.Errorf("parsing body: %v", err)
	}

	// Проверяем, что пользователь с таким email еще не зарегистрирован
	exists := h.storage.Exists(regReq.Email)
	if exists {
		return errors.New("the user already exists")
	}

	// Сохраняем нового зарегистрированного пользователя
	h.storage.Add(regReq.Email, &models.User{
		Email:    regReq.Email,
		Name:     regReq.Name,
		Password: regReq.Password,
	})

	return c.SendStatus(fiber.StatusCreated)
}

// Обработчик HTTP-запросов на вход в аккаунт
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	regReq := LoginRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		return fmt.Errorf("parsing body: %v", err)
	}

	// Ищем пользователя по электронной почте
	user, err := h.storage.Get(regReq.Email)
	// Если пользователь не найден, возвращаем ошибку
	if err != nil {
		return errBadCredentials
	}
	// Если пользователь найден, но у него другой пароль, возвращаем ошибку
	if user.Password != regReq.Password {
		return errBadCredentials
	}

	// Генерируем JWT-токен для пользователя,
	// который он будет использовать в будущих HTTP-запросах

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(h.signKey)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(LoginResponse{AccessToken: t})
}
