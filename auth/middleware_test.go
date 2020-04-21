package auth

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/service"
	"github.com/erickoliv/finances-api/test/entities"
	"github.com/erickoliv/finances-api/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMiddleware(t *testing.T) {

	validCookie := http.Cookie{
		Name:  domain.AuthCookie,
		Value: "a value",
	}

	tests := []struct {
		name     string
		signer   func() service.Signer
		cookie   http.Cookie
		status   int
		response string
	}{
		{
			name: "validate request without auth cookie",
			signer: func() service.Signer {
				return &mocks.Signer{}
			},
			status:   http.StatusUnauthorized,
			response: `{"message":"auth cookie missing"}`,
		},
		{
			name: "error to validate auth cookie",
			signer: func() service.Signer {
				signer := &mocks.Signer{}

				signer.On("Validate", mock.Anything, validCookie.Value).Return(uuid.Nil, errors.New("invalid auth token"))
				return signer
			},
			status:   http.StatusUnauthorized,
			cookie:   validCookie,
			response: `{"message":"invalid auth token"}`,
		},
		{
			name: "successfully validate auth token and set user uuid into context",
			signer: func() service.Signer {
				signer := &mocks.Signer{}

				signer.On("Validate", mock.Anything, validCookie.Value).Return(entities.ValidUser().UUID, nil)
				return signer
			},
			status:   http.StatusOK,
			cookie:   validCookie,
			response: entities.ValidUser().UUID.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			got := Middleware(tt.signer())

			router.Use(got)
			router.GET("", func(c *gin.Context) {
				user := c.MustGet(domain.LoggedUser).(uuid.UUID)

				c.JSON(http.StatusOK, user.String())
			})

			req, _ := http.NewRequest("GET", "", nil)
			req.AddCookie(&tt.cookie)
			// req.Header.Add("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}
