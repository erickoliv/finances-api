package api

import (
	"os"
	"testing"

	"github.com/google/uuid"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var router = gin.Default()
var user = uuid.New()

func TestMain(m *testing.M) {
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true
	// GORM
	db, _ := gorm.Open(mocket.DriverName, "test") // Can be any connection string
	defer db.Close()

	router.Use(database.Middleware(db))
	router.Use(dummyUser())

	setupTagsDatabase()

	Routes(router.Group("api"))
	os.Exit(m.Run())
}

func dummyUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Set(common.LoggedUser, user)
		print(user.String())
		c.Next()
	}
}
