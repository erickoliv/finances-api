package db

import (
	"fmt"
	"log"
	"os"

	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/categories"
	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/entries"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // for pg dialect
)

// Prepare database connection
func Prepare() *gorm.DB {

	host := getEnvConfig("DB_HOST")
	user := getEnvConfig("DB_USER")
	password := getEnvConfig("DB_PASSWORD")
	port := getEnvConfig("DB_PORT")
	database := getEnvConfig("DB_NAME")

	dbURL := fmt.Sprintf("host=%s user=%s port=%s dbname=%s password=%s sslmode=disable", host, user, port, database, password)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	// database migrations, pendind real migration startup using plain sql scripts
	db.AutoMigrate(&domain.Tag{})
	db.AutoMigrate(&auth.User{})
	db.AutoMigrate(&categories.Category{})
	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&entries.Entry{})
	db.AutoMigrate(&domain.EntryTag{})

	return db
}

func getEnvConfig(s string) string {
	if value, found := os.LookupEnv(s); found {
		return value
	}

	log.Fatalf("Environment variable %s not found", s)
	return ""
}

// Middleware adds a gorm.DB connection pool reference inside gin.Context
func Middleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(domain.DB, db)
		c.Next()
	}
}
