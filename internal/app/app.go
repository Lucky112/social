package app

import (
	"context"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/service"
	"github.com/Lucky112/social/internal/transport"
)

func Run(ctx context.Context, config *config.Config) {
	service, err := service.NewService(ctx, config.StorageConfig)
	if err != nil {
		panic(err)
	}
	defer service.Close(ctx)

	server := transport.NewServer(
		config.ServerConfig,
		service.AuthService(),
		service.ProfilesService(),
		service.FriendsService(),
		service.PostsService(),
	)

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
