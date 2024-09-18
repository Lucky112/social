package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBConfig     *DBConfig     `json:"db_config"     yaml:"db_config"     validate:"required"`
	ServerConfig *ServerConfig `json:"server_config" yaml:"server_config" validate:"required"`
}

type DBConfig struct {
	Host     string `json:"host"     yaml:"host"     validate:"required,hostname|ip"`
	Port     uint16 `json:"port"     yaml:"port"     validate:"required,min=1,max=65535"`
	User     string `json:"user"     yaml:"user"     validate:"required"`
	Password string `json:"password" yaml:"password" validate:"required"`
	Database string `json:"database" yaml:"database" validate:"required"`
}

type ServerConfig struct {
	Port   uint16 `json:"port" yaml:"port" validate:"required,min=1,max=65535"`
	JWTKey string `json:"jwt_key" yaml:"jwt_key" validate:"required"`
}

func Load(filename string) (*Config, error) {
	var config *Config
	var err error

	ext := filepath.Ext(filename)

	switch ext {
	case ".json":
		config, err = fromFile(filename, json.Unmarshal)
		if err != nil {
			return nil, fmt.Errorf("loading from json: %v", err)
		}
	case ".yaml", ".yml":
		config, err = fromFile(filename, yaml.Unmarshal)
		if err != nil {
			return nil, fmt.Errorf("loading from yaml: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file type: .json or .yaml expected")
	}

	err = validate(config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	return config, nil
}

func fromFile(filename string, parser func([]byte, any) error) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}

	var config Config
	err = parser(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %v", err)
	}

	return &config, nil
}

func validate(config *Config) error {
	validate := validator.New()
	err := validate.Struct(config)
	return err
}
