package postgres

import "fmt"

type Config struct {
	User     string
	Password string
	Database string
	Host     string
	Port     uint16
}

func (cfg Config) connectionURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
