package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestExtractUser(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name           string
		want           uuid.UUID
		prepareContext func(c *gin.Context)
		err            error
	}{
		{
			name: "successfully gets a valid user uuid from context using extractor",
			want: validUUID,
			prepareContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validUUID.String())
			},
			err: nil,
		},
		{
			name: "error when the context doesn't have a user uuid",
			want: uuid.Nil,
			prepareContext: func(c *gin.Context) {
			},
			err: errors.New("user not present in context"),
		},
		{
			name: "error when the context contains a invalid data into logged user constant",
			want: uuid.Nil,
			prepareContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, "invalid value")
			},
			err: errors.New("invalid UUID length: 13"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			router.Use(tt.prepareContext)

			router.GET("", func(c *gin.Context) {
				user, err := ExtractUser(c)

				assert.Equal(t, tt.want, user)
				assert.Equal(t, tt.err, err)

				c.JSON(http.StatusOK, ":)")
			})

			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		})
	}
}

func TestExtractUUID(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name string
		url  string
		want uuid.UUID
		err  error
	}{
		{
			name: "successfully gets a uuid from url param using extractor",
			want: validUUID,
			url:  fmt.Sprintf("/%s", validUUID.String()),
			err:  nil,
		},
		{
			name: "error to get a valid uuid using extractor",
			url:  fmt.Sprintf("/invalid-value"),
			err:  errors.New("invalid uuid in context: invalid UUID length: 13"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/:uuid", func(c *gin.Context) {
				identifier, err := ExtractUUID(c)

				assert.Equal(t, tt.want, identifier)
				assert.Equal(t, tt.err, err)

				c.JSON(http.StatusOK, ":)")
			})

			req, _ := http.NewRequest("GET", tt.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		})
	}
}
