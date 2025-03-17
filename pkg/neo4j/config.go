package neo4j

import "fmt"

type Config struct {
	User     string
	Password string
	Host     string
	Port     uint16
}

func (cfg Config) connectionURL() string {
	return fmt.Sprintf("bolt://%s:%d/", cfg.Host, cfg.Port)
}
