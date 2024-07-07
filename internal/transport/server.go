package transport

import (
	"fmt"

	"github.com/Lucky112/social/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	server *fiber.App
}

const signKey = "encription-key"

func NewServer() Server {
	authStorage := storage.NewAuthStorage()
	authHandler := NewAuthHandler(authStorage, signKey)

	server := fiber.New()

	server.Use(recover.New())

	publicGroup := server.Group("")
	publicGroup.Post("/register", authHandler.Register)
	publicGroup.Post("/login", authHandler.Login)

	return Server{
		server: server,
	}
}

func (s Server) Start(port int16) error {
	address := fmt.Sprintf(":%d", port)

	err := s.server.Listen(address)
	return err
}
