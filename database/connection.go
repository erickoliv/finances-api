package database

import (
	"log"
	"os"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // for pg dialect
)

// PrepareDatabase is the function designed
// to prepare and return all shared/singleton application props
func PrepareDatabase() *gorm.DB {
	dbURL := getEnvConfig("DB_URL")
	_ = getEnvConfig("APP_TOKEN")

	log.Println("url", dbURL)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	// database migrations
	db.AutoMigrate(&model.Tag{})
	db.AutoMigrate(&model.User{})

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
		c.Set(common.DB, db)
		c.Next()
	}
}
