package transport

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/transport/auth"
	"github.com/Lucky112/social/internal/transport/jwt"
	"github.com/Lucky112/social/internal/transport/profiles"
)

type Server struct {
	server *fiber.App
	port   uint16
}

func NewServer(cfg *config.ServerConfig, authService auth.AuthService, profilesService profiles.ProfilesService) Server {
	jwtKey := []byte(cfg.JWTKey)

	authHandler := auth.NewAuthHandler(authService, jwtKey)
	profilesHandler := profiles.NewProfilesHandler(profilesService)

	server := fiber.New()

	server.Use(recover.New())

	publicGroup := server.Group("")
	publicGroup.Post("/register", authHandler.Register)
	publicGroup.Post("/login", authHandler.Login)

	authorizedGroup := server.Group("")
	authorizedGroup.Use(jwt.Middleware(jwtKey))

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
