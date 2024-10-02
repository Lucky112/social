package service

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/config"
	pg "github.com/Lucky112/social/internal/storage/postgres"
	"github.com/Lucky112/social/pkg/postgres"
)

type Service struct {
	dbpool postgres.Pool
}

func NewService(ctx context.Context, config *config.DBConfig) (Service, error) {
	cfg := toPostgresConfig(config)

	err := migrateDB(cfg)
	if err != nil {
		return Service{}, fmt.Errorf("migrating database: %v", err)
	}

	dbpool, err := postgres.ViaPGX(ctx, cfg)
	if err != nil {
		return Service{}, fmt.Errorf("creating pgx pool: %v", err)
	}

	return Service{
		dbpool: dbpool,
	}, err
}

func (s Service) AuthService() AuthService {
	storage := pg.NewUsersProvider(s.dbpool)
	return NewAuthService(storage)
}

func (s Service) ProfilesService() ProfilesService {
	storage := pg.NewProfilesProvider(s.dbpool)
	return NewProfilesService(storage)
}

func toPostgresConfig(cfg *config.DBConfig) *postgres.Config {
	return &postgres.Config{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		Host:     cfg.Host,
		Port:     cfg.Port,
	}
}

func migrateDB(cfg *postgres.Config) error {
	sqldb, err := postgres.ViaSTD(cfg)
	if err != nil {
		return fmt.Errorf("opening db: %v", err)
	}

	err = pg.ApplyMigrations(sqldb)
	if err != nil {
		return fmt.Errorf("applying db migrations: %v", err)
	}

	err = sqldb.Close()
	if err != nil {
		return fmt.Errorf("closing db after migration: %v", err)
	}

	return nil
}
