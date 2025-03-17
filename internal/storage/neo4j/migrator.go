package neo4j

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migrator "github.com/golang-migrate/migrate/v4/database/neo4j"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//go:embed migrations/*.cypher
var fs embed.FS

func ApplyMigrations(driver neo4j.Driver) error {
	inst, err := migrator.WithInstance(driver, &migrator.Config{})
	if err != nil {
		return fmt.Errorf("creating migrator client: %v", err)
	}

	files, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", files, "neo4j", inst)
	if err != nil {
		return fmt.Errorf("creating migrator: %v", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// it's ok - already migrated
			return nil
		}

		return fmt.Errorf("migrating up: %v", err)
	}

	return nil
}

func RollbackMigrations(driver neo4j.Driver) error {
	inst, err := migrator.WithInstance(driver, &migrator.Config{})
	if err != nil {
		return fmt.Errorf("creating migrator client: %v", err)
	}

	files, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", files, "neo4j", inst)
	if err != nil {
		return fmt.Errorf("creating migrator: %v", err)
	}

	err = m.Down()
	if err != nil {
		return fmt.Errorf("migrating down: %v", err)
	}

	return nil
}
