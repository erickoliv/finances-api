package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var router = gin.Default()

func TestIndexHandler(t *testing.T) {
	// router = gin.New()
	router.GET("/", IndexHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"status\":\"ok\"")
}

func TestMain(m *testing.M) {
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true
	// GORM
	db, _ := gorm.Open(mocket.DriverName, "connection_string") // Can be any connection string
	defer db.Close()
	
	router.Use(database.Middleware(db))
	TagRoutes(router)
	os.Exit(m.Run())

}
