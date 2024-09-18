package config

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
}

type DBConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}

type ServerConfig struct {
	Port uint16
}
