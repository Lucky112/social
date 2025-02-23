package service

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/storage/inmemory"
	n4j "github.com/Lucky112/social/internal/storage/neo4j"
	pg "github.com/Lucky112/social/internal/storage/postgres"
	"github.com/Lucky112/social/pkg/neo4j"
	"github.com/Lucky112/social/pkg/postgres"
)

type Service struct {
	dbpool      *postgres.Pool
	neo4jDriver *neo4j.Driver
}

func NewService(ctx context.Context, config *config.StorageConfig) (Service, error) {
	dbpool, err := connectToPostgres(ctx, config.DBConfig)
	if err != nil {
		return Service{}, fmt.Errorf("connecting to postgres: %v", err)
	}

	neo4jDriver, err := connectToNeo4j(ctx, config.Neo4jConfig)
	if err != nil {
		return Service{}, fmt.Errorf("connecting to neo4j: %v", err)
	}

	return Service{
		dbpool:      dbpool,
		neo4jDriver: neo4jDriver,
	}, err
}

func (s Service) Close(ctx context.Context) error {
	s.dbpool.Close()
	return s.neo4jDriver.Close(ctx)
}

func connectToPostgres(ctx context.Context, config *config.DBConfig) (*postgres.Pool, error) {
	cfg := toPostgresConfig(config)

	err := migrateDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("migrating database: %v", err)
	}

	dbpool, err := postgres.ViaPGX(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating pgx pool: %v", err)
	}

	return &dbpool, nil
}

func connectToNeo4j(ctx context.Context, config *config.Neo4jConfig) (*neo4j.Driver, error) {
	cfg := toNeo4jConfig(config)

	// err := migrateNeo4j(cfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("migrating database: %v", err)
	// }

	driver, err := neo4j.NewDriver(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating neo4j driver: %v", err)
	}

	return driver, nil
}

func (s Service) AuthService() AuthService {
	storage := pg.NewUsersProvider(s.dbpool)
	return NewAuthService(storage)
}

func (s Service) ProfilesService() ProfilesService {
	storage := pg.NewProfilesProvider(s.dbpool)
	return NewProfilesService(storage)
}

func (s Service) FriendsService() FriendsService {
	// storage := n4j.NewFriendsProvider(s.neo4jDriver)
	storage := inmemory.NewFriendsStorage()
	return NewFriendsService(storage)
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

func toNeo4jConfig(cfg *config.Neo4jConfig) *neo4j.Config {
	return &neo4j.Config{
		User:     cfg.User,
		Password: cfg.Password,
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

func migrateNeo4j(cfg *neo4j.Config) error {
	driver, err := neo4j.NewStdDriver(cfg)
	if err != nil {
		return fmt.Errorf("opening neo4j: %v", err)
	}

	err = n4j.ApplyMigrations(driver)
	if err != nil {
		return fmt.Errorf("applying neo4j migrations: %v", err)
	}

	err = driver.Close()
	if err != nil {
		return fmt.Errorf("closing neo4j after migration: %v", err)
	}

	return nil
}
