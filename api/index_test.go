package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
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
	var db = database.PrepareDatabase()
	defer db.Close()
	router.Use(database.Middleware(db))
	TagRoutes(router)
	os.Exit(m.Run())
}
