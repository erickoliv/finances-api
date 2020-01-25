package rest

// import (
// 	"os"
// 	"testing"

// 	"github.com/google/uuid"

// 	mocket "github.com/Selvatico/go-mocket"
// 	"github.com/erickoliv/finances-api/domain"
// 	"github.com/gin-gonic/gin"
// 	"github.com/jinzhu/gorm"
// )

// var router = gin.Default()
// var user = uuid.New()

// func TestMain(m *testing.M) {
// 	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
// 	mocket.Catcher.Logging = true
// 	// GORM
// 	db, _ := gorm.Open(mocket.DriverName, "test") // Can be any connection string
// 	defer db.Close()

// 	router.Use(dummyMiddleware(db))

// 	setupTagsDatabase()

// 	// Routes(router.Group("api"))
// 	os.Exit(m.Run())
// }

// func dummyMiddleware(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		c.Set(domain.DB, db)
// 		c.Set(domain.LoggedUser, user)
// 		print(user.String())
// 		c.Next()
// 	}
// }
