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

func TestList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("Status Diferente de 404")
	}

}

func TestMain(m *testing.M) {
	var db = database.PrepareDatabase()
	defer db.Close()

	TagRoutes(router)

	defer db.Close()

	os.Exit(m.Run())
}

func TestTagRoutes(t *testing.T) {
	tests := []struct {
		name string
		r    *gin.Engine
	}{
		{"simple test", router},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TagRoutes(tt.r)
		})
	}
}

func TestGetTags(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetTags(tt.args.c)
		})
	}
}
