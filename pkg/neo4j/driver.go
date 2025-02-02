package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Driver struct {
	driver neo4j.DriverWithContext
}

func NewDriver(ctx context.Context, config *Config) (*Driver, error) {
	driver, err := neo4j.NewDriverWithContext(config.URI, neo4j.BasicAuth(config.User, config.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("creating Neo4j driver: %w", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to Neo4j: %w", err)
	}

	return &Driver{driver: driver}, nil
}

func (d *Driver) NewReadSession(ctx context.Context) neo4j.SessionWithContext {
	return d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
}

func (d *Driver) NewWriteSession(ctx context.Context) neo4j.SessionWithContext {
	return d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
}

func (d *Driver) Close(ctx context.Context) error {
	return d.driver.Close(ctx)
}
