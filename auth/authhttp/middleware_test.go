package authhttp

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/auth/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {

	validCookie := http.Cookie{
		Name:  auth.AuthCookie,
		Value: "a value",
	}

	tests := []struct {
		name     string
		signer   func() auth.SessionSigner
		cookie   http.Cookie
		status   int
		response string
	}{
		{
			name: "validate request without auth cookie",
			signer: func() auth.SessionSigner {
				return &mocks.SessionSigner{}
			},
			status:   http.StatusUnauthorized,
			response: `{"message":"auth cookie missing"}`,
		},
		{
			name: "error to validate auth cookie",
			signer: func() auth.SessionSigner {
				signer := &mocks.SessionSigner{}

				signer.On("Validate", validCookie.Value).Return("", errors.New("invalid auth token"))
				return signer
			},
			status:   http.StatusUnauthorized,
			cookie:   validCookie,
			response: `{"message":"invalid auth token"}`,
		},
		{
			name: "successfully validate auth token and set user uuid into context",
			signer: func() auth.SessionSigner {
				signer := &mocks.SessionSigner{}

				signer.On("Validate", validCookie.Value).Return(mocks.ValidUser().UUID.String(), nil)
				return signer
			},
			status:   http.StatusOK,
			cookie:   validCookie,
			response: mocks.ValidUser().UUID.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			got := Middleware(tt.signer())

			router.Use(got)
			router.GET("", func(c *gin.Context) {
				user := c.MustGet(auth.LoggedUser).(string)

				c.JSON(http.StatusOK, user)
			})

			req, _ := http.NewRequest("GET", "/", nil)
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
