package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateTag(t *testing.T) {
	text := "{\"name\":\"nova tag\",\"description\": \"descrição com cacteres espéciàis\"}"

	req, _ := http.NewRequest("POST", "/api/tags", strings.NewReader(text))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
