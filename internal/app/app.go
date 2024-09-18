package app

import (
	"context"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/service"
	"github.com/Lucky112/social/internal/transport"
)

func Run() {
	config := config.Config{}

	service, err := service.NewService(context.Background(), config.DBConfig)
	if err != nil {
		panic(err)
	}

	server := transport.NewServer(config.ServerConfig, service.AuthService(), service.ProfilesService())

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
