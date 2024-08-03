package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

type Config struct {
	MinLogLevel string `yaml:"min_log_level" validate:"required,oneof=debug info warn error"`
	PostgreSQL  struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required,numeric"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Name     string `yaml:"database" validate:"required"`
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

	envConfigFilePath := os.Getenv("ENV_CONFIG_FILE_PATH")
	if err = godotenv.Load(envConfigFilePath); err != nil {
		log.Fatal("Error loading .env file")
	}

	// postgresql
	config.PostgreSQL.Host = os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	config.PostgreSQL.Port, err = strconv.Atoi(dbPort)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	config.PostgreSQL.User = os.Getenv("DB_USER")
	config.PostgreSQL.Password = os.Getenv("DB_PASSWORD")
	config.PostgreSQL.Name = os.Getenv("DB_NAME")

	// redis
	config.Redis.Host = os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	config.Redis.Port, err = strconv.Atoi(redisPort)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	config.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	config.JwtSecret = os.Getenv("JWT_SECRET")

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
