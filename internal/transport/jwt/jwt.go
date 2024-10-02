package jwt

import (
	"fmt"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const jwtContextKey = "user"
const userIdClaim = "sub"

func MakeToken(userId string, key []byte) (string, error) {
	payload := jwt.MapClaims{
		userIdClaim: userId,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("signing token: %v", err)
	}

	return t, nil
}

func Middleware(key []byte) any {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: key,
		},
		ContextKey: jwtContextKey,
	})
}

func ExtractUserId(c *fiber.Ctx) (string, error) {
	user, ok := c.Locals(jwtContextKey).(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("no jwt token found in '%s' local", jwtContextKey)
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unable to convert %v to jwt.MapClaims", user.Claims)
	}

	userId, ok := claims[userIdClaim].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract userId string from %v", claims[userIdClaim])
	}

	return userId, nil
}
