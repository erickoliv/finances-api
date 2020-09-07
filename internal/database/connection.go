package database

import (
	"fmt"

	"github.com/erickoliv/finances-api/internal/cfg"
	"gorm.io/driver/postgres" // for pg dialect
	"gorm.io/gorm"
)

const connectionString = "host=%s user=%s port=%s dbname=%s password=%s sslmode=disable"

// Connect database connection
func Connect(config *cfg.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(connectionString, config.Host, config.User, config.Port, config.Schema, config.Password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
