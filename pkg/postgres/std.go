package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
)

func ViaSTD(config *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.connectionURL())
	if err != nil {
		return nil, fmt.Errorf("opening database: %v", err)
	}

	return db, nil
}
