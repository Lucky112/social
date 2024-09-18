package transport

import (
	"fmt"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/transport/auth"
	"github.com/Lucky112/social/internal/transport/profiles"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	server *fiber.App
	port   uint16
}

const signKey = "encription-key"

func NewServer(cfg config.ServerConfig, authService auth.AuthService, profilesService profiles.ProfilesService) Server {
	authHandler := auth.NewAuthHandler(authService, signKey)
	profilesHandler := profiles.NewProfilesHandler(profilesService)

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
		port:   cfg.Port,
	}
}

func (s Server) Start() error {
	address := fmt.Sprintf(":%d", s.port)

	err := s.server.Listen(address)
	return err
}
