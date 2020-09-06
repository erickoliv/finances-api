package cfg

import "time"

type (
	Database struct {
		Host          string
		User          string
		Password      string
		Port          string
		Schema        string
		EnableLogging bool
	}

	Auth struct {
		Token string
		TTL   time.Duration
	}

	Config struct {
		JWT Auth
		DB  Database
	}

	envFetcher func(string) string
)

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

	return &config, nil
}
