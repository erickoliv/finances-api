package database

import (
	"log"
	"os"

	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	// database migrations
	db.AutoMigrate(&model.Tag{})

	return db
}

func getEnvConfig(s string) string {
	if value, found := os.LookupEnv(s); found {
		return value
	} else {
		log.Fatalf("Environment variable %s not found", s)
	}
	return ""
}
