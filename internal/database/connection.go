package database

import (
	"fmt"

	"github.com/erickoliv/finances-api/internal/cfg"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // for pg dialect
)

const connectionString = "host=%s user=%s port=%s dbname=%s password=%s sslmode=disable"

// Connect database connection
func Connect(config *cfg.Database) (*gorm.DB, error) {
	dbURL := fmt.Sprintf(connectionString, config.Host, config.User, config.Port, config.Schema, config.Password)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.LogMode(config.EnableLogging)

	return db, nil
}
