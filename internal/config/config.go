package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	HttpPort     int    `yaml:"http_port" validate:"required,numeric"`
	ReadTimeout  int    `yaml:"read_timeout" validate:"required,numeric"`
	WriteTimeout int    `yaml:"write_timeout" validate:"required,numeric"`
	MinLogLevel  string `yaml:"min_log_level" validate:"required,oneof=debug info warn error"`
	PostgreSQL   struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required,numeric"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Name     string `yaml:"name" validate:"required"`
	} `yaml:"postgresql"`
	Redis struct {
		Host string `yaml:"host" validate:"required"`
		Port int    `yaml:"port" validate:"required,numeric"`
	} `yaml:"redis"`
	Telegram struct {
		Token string `yaml:"token" validate:"required"`
	} `yaml:"telegram"`
	JwtSecret string `yaml:"jwt_secret" validate:"required"`
}

func LoadConfig() (*Config, error) {
	yamlConfigFilePath := os.Getenv("YAML_CONFIG_FILE_PATH")
	if yamlConfigFilePath == "" {
		return nil, fmt.Errorf("yaml config file path is not set")
	}

	f, err := os.Open(yamlConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Printf("unable to close config file: %v", err)
		}
	}(f)

	var config Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
