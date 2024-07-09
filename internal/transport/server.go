package transport

import (
	"fmt"

	"github.com/Lucky112/social/internal/storage"
	"github.com/Lucky112/social/internal/transport/auth"
	"github.com/Lucky112/social/internal/transport/profiles"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	server *fiber.App
}

const signKey = "encription-key"

func NewServer() Server {
	authStorage := storage.NewAuthStorage()
	authHandler := auth.NewAuthHandler(authStorage, signKey)

	profileStorage := storage.NewProfileStorage()
	profilesHandler := profiles.NewProfilesHandler(profileStorage)

	server := fiber.New()

	server.Use(recover.New())

	publicGroup := server.Group("")
	publicGroup.Post("/register", authHandler.Register)
	publicGroup.Post("/login", authHandler.Login)

	authorizedGroup := server.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(signKey),
		},
	}))

	authorizedGroup.Post("/profiles", profilesHandler.CreateProfile)
	authorizedGroup.Get("/profiles", profilesHandler.GetProfiles)
	authorizedGroup.Get("/profiles/:id", profilesHandler.GetProfile)

	return Server{
		server: server,
	}
}

func (s Server) Start(port int16) error {
	address := fmt.Sprintf(":%d", port)

	err := s.server.Listen(address)
	return err
}
