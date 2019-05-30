package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ericktm/olivsoft-golang-api/database"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func TestMain(m *testing.M) {
	var db = database.PrepareDatabase()
	defer db.Close()
	router.Use(database.Middleware(db))
	TagRoutes(router)

	os.Exit(m.Run())
}

func TestGetTags(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status %d; want 200", w.Code)
	}
}
