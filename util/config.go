package util

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ericktm/olivsoft-golang-api/api"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Config is the Application Shared/Singleton props resource
type Config struct {
	ApplicationName string
	DB              *gorm.DB
	Router          *gin.Engine
	StartupTime     time.Time
}

// GetConfig is the function designed
// to prepare and return all shared/singleton application props
func GetConfig() Config {
	dbURL := getEnvConfig("DB_URL")
	_ = getEnvConfig("APP_TOKEN")

	log.Println("url", dbURL)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(true)

	// database migrations
	log.Println("start database migrations")
	db.AutoMigrate(&model.Tag{})
	log.Println("stop database migrations")

	r := gin.Default()
	r.Use(authMiddleware())

	r.GET("/", api.IndexHandler(db))
	r.GET("/api/tags/:uuid", api.GetTag(db))
	r.PUT("/api/tags/:uuid", api.UpdateTag(db))
	r.DELETE("/api/tags/:uuid", api.DeleteTag(db))
	r.GET("/api/tags", api.GetTags(db))
	r.POST("/api/tags", api.CreateTag(db))
	//404 handler
	// r.NotFoundHandler = api.PageNotFound()
	cfg := Config{
		"OlivSoft",
		db,
		r,
		time.Now(),
	}

	return cfg
}

// dummy token auth
func authMiddleware() gin.HandlerFunc {
	envToken := os.Getenv("APP_TOKEN")

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, api.ErrorMessage{Message: "missing authentication token"})
			c.Abort()
			return
		}

		if token != envToken {
			c.JSON(http.StatusUnauthorized, api.ErrorMessage{Message: "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getEnvConfig(s string) string {
	if value, found := os.LookupEnv(s); found {
		return value
	} else {
		log.Fatalf("Environment variable %s not found", s)
	}
	return ""
}
