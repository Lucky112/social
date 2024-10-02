package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/Lucky112/social/pkg/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

var db *sql.DB

func TestMain(m *testing.M) {
	const (
		dbname   = "demodb"
		user     = "user"
		password = "pwd"
	)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	envs := []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		fmt.Sprintf("POSTGRES_USER=%s", user),
		fmt.Sprintf("POSTGRES_DB=%s", dbname),
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "16.3", envs)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(600)

	hostport := resource.GetHostPort("5432/tcp")
	host, portstr, err := net.SplitHostPort(hostport)
	if err != nil {
		log.Fatalf("Could not extraxt host and port from %s", hostport)
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatalf("Could not convert port '%s' to integer", portstr)
	}

	config := postgres.Config{
		User:     user,
		Password: password,
		Database: dbname,
		Host:     host,
		Port:     uint16(port),
	}

	db, err = postgres.ViaSTD(&config)
	if err != nil {
		log.Fatalf("Could not connect to db: %s", err)
	}
	for i := 0; i < 60; i++ {
		if db.Ping() == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
	defer func() {
		err := pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}

	}()

	m.Run()
}

func TestMigrations(t *testing.T) {
	err := ApplyMigrations(db)
	require.NoError(t, err)

	err = RollbackMigrations(db)
	require.NoError(t, err)

	err = ApplyMigrations(db)
	require.NoError(t, err)
}
