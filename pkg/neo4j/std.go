package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func NewStdDriver(config *Config) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(config.URI, neo4j.BasicAuth(config.User, config.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("creating Neo4j driver: %w", err)
	}

	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, fmt.Errorf("connecting to Neo4j: %w", err)
	}

	return driver, nil
}
