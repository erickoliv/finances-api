package cfg

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type (
	Database struct {
		Host          string `validate:"required"`
		User          string `validate:"required"`
		Password      string `validate:"required"`
		Port          string `validate:"required,numeric"`
		Schema        string `validate:"required"`
		EnableLogging bool
	}

	Auth struct {
		Token string        `validate:"required"`
		TTL   time.Duration `validate:"required"`
	}

	Config struct {
		JWT Auth     `validate:"required"`
		DB  Database `validate:"required"`
	}

	envFetcher func(string) string
)

var validate = validator.New()

func Load(getEnv envFetcher) (*Config, error) {
	config := Config{
		JWT: Auth{
			Token: getEnv("APP_TOKEN"),
			TTL:   time.Hour,
		},
		DB: Database{
			Host:          getEnv("DB_HOST"),
			User:          getEnv("DB_USER"),
			Password:      getEnv("DB_PASSWORD"),
			Port:          getEnv("DB_PORT"),
			Schema:        getEnv("DB_NAME"),
			EnableLogging: true,
		},
	}
	err := validate.Struct(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
