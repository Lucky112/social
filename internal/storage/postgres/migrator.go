package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func applyMigrations(db *sql.DB) error {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("creating migrator db client: %v", err)
	}

	files, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", files, "postgres", driver)
	if err != nil {
		return fmt.Errorf("creating migrator: %v", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("migrating up: %v", err)
	}

	return nil
}

func rollbackMigrations(db *sql.DB) error {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("creating migrator db client: %v", err)
	}

	files, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("creating iofs driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", files, "postgres", driver)
	if err != nil {
		return fmt.Errorf("creating migrator: %v", err)
	}

	err = m.Down()
	if err != nil {
		return fmt.Errorf("migrating down: %v", err)
	}

	return nil
}
