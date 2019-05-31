package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {

	tag := model.Tag{
		Name:        "a simple tag",
		Description: "a tag description",
	}

	str, _ := json.Marshal(tag)

	req, _ := http.NewRequest("POST", "/api/tags", bytes.NewReader(str))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetTags(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	println(w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}
