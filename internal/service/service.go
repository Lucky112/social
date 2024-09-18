package service

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/config"
	pg "github.com/Lucky112/social/internal/storage/postgres"
	"github.com/Lucky112/social/pkg/postgres"
)

type Service struct {
	db postgres.Pool
}

func NewService(ctx context.Context, config config.DBConfig) (Service, error) {
	cfg := toPostgresConfig(config)

	db, err := postgres.ViaPGX(ctx, cfg)
	if err != nil {
		return Service{}, fmt.Errorf("creating db: %v", err)
	}

	return Service{
		db: db,
	}, err
}

func (s Service) AuthService() AuthService {
	storage := pg.NewUsersProvider(s.db)
	return NewAuthService(storage)
}

func (s Service) ProfilesService() ProfilesService {
	storage := pg.NewProfilesProvider(s.db)
	return NewProfilesService(storage)
}

func toPostgresConfig(cfg config.DBConfig) postgres.Config {
	return postgres.Config{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		Host:     cfg.Host,
		Port:     cfg.Port,
	}
}
