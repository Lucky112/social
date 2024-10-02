package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := Config{
		User:     "A",
		Password: "123",
		Database: "db",
		Host:     "host",
		Port:     10000,
	}

	expectedConnectionURL := "postgres://A:123@host:10000/db"
	actual := cfg.connectionURL()

	assert.Equal(t, expectedConnectionURL, actual)
}

func TestNewPool(t *testing.T) {
	t.Run("Correct config", func(t *testing.T) {
		cfg := Config{
			User:     "A",
			Password: "123",
			Database: "db",
			Host:     "host",
			Port:     10000,
		}

		_, err := ViaPGX(context.Background(), &cfg)

		assert.NoError(t, err)
	})

	t.Run("Empty config", func(t *testing.T) {
		cfg := Config{}

		_, err := ViaPGX(context.Background(), &cfg)

		assert.Error(t, err)
	})
}

func TestStd(t *testing.T) {
	t.Run("Correct config", func(t *testing.T) {
		cfg := Config{
			User:     "A",
			Password: "123",
			Database: "db",
			Host:     "host",
			Port:     10000,
		}

		_, err := ViaSTD(&cfg)

		assert.NoError(t, err)
	})
}
