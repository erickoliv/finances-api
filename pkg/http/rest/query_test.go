package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExtractFilters(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name           string
		needUser       bool
		url            string
		prepareContext func(c *gin.Context)
		want           *Query
		err            error
	}{
		{
			name:           "error to get user from context when needed",
			needUser:       true,
			url:            "/resource",
			prepareContext: func(*gin.Context) {},
			want:           nil,
			err:            errors.New("user not present in context"),
		},
		{
			name:     "successfuly get all posible filters",
			needUser: true,
			url:      "/resource?limit=10&page=3&sort=age&q_name__eq=name&q_city__like=califo&q_age__gte=18&q_age__lte=60&ignorewithout_q=yes",
			prepareContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validUUID.String())
				c.Next()
			},
			want: &Query{
				Page:  3,
				Limit: 10,
				Sort:  "age",
				Filters: map[string]interface{}{
					"owner = ?":   validUUID,
					"name = ?":    "name",
					"city LIKE ?": "%califo%",
					"age >= ?":    "18",
					"age <= ?":    "60",
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(tt.prepareContext)

			router.GET("/resource", func(c *gin.Context) {
				query, err := ExtractFilters(c, tt.needUser)

				assert.Equal(t, tt.want, query)
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
