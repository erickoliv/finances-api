package api

import (
	"os"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var router = gin.Default()

func TestMain(m *testing.M) {
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true
	// GORM
	db, _ := gorm.Open(mocket.DriverName, "connection_string") // Can be any connection string
	defer db.Close()

	router.Use(database.Middleware(db))
	setupTagsDatabase()

	Routes(router.Group("api"))
	os.Exit(m.Run())
}
