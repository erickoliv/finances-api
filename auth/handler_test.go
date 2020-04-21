package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/erickoliv/finances-api/service"
	"github.com/erickoliv/finances-api/test/entities"
	"github.com/erickoliv/finances-api/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestHandler() HTTPHandler {
	authenticator := &mocks.Authenticator{}
	signer := &mocks.Signer{}

	return NewHTTPHandler(authenticator, signer)
}

func TestNewHTTPHandler(t *testing.T) {
	authenticator := &mocks.Authenticator{}
	signer := &mocks.Signer{}

	tests := []struct {
		name          string
		authenticator service.Authenticator
		signer        service.Signer
		want          HTTPHandler
	}{
		{
			name:          "creates a new http handler",
			authenticator: authenticator,
			signer:        signer,
			want: &httpHandler{
				auth: authenticator,
				sign: signer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := NewHTTPHandler(tt.authenticator, tt.signer)

			assert.Equal(t, tt.want, view)

			router := gin.New()
			group := router.Group("")
			view.Router(group)
		})
	}
}

func Test_httpHandler_login(t *testing.T) {
	validCredential := credential{
		Username: "user",
		Password: "userpass",
	}

	tests := []struct {
		name         string
		setupHandler func() HTTPHandler
		setupContext func(c *gin.Context)
		body         string
		status       int
		response     string
	}{
		{
			name: "successfully login a valid user",
			setupHandler: func() HTTPHandler {
				authenticator := &mocks.Authenticator{}
				signer := &mocks.Signer{}

				authenticator.On("Login", mock.Anything, validCredential.Username, validCredential.Password).Return(entities.ValidUser(), nil)
				signer.On("SignUser", mock.Anything, entities.ValidUser()).Return("token", nil)

				return NewHTTPHandler(authenticator, signer)
			},
			body:     serialize(&validCredential),
			status:   http.StatusOK,
			response: `{"message":"token"}`,
		},
		{
			name: "error to login due to invalid payload data",
			setupHandler: func() HTTPHandler {
				authenticator := &mocks.Authenticator{}
				signer := &mocks.Signer{}
				return NewHTTPHandler(authenticator, signer)
			},
			body:     "{}",
			status:   http.StatusBadRequest,
			response: `{"message":"invalid payload"}`,
		},
		{
			name: "error to login due to invalid credentials",
			setupHandler: func() HTTPHandler {
				authenticator := &mocks.Authenticator{}
				signer := &mocks.Signer{}

				authenticator.On("Login", mock.Anything, validCredential.Username, validCredential.Password).Return(nil, errors.New("invalid credentials"))
				return NewHTTPHandler(authenticator, signer)
			},
			body:     serialize(&validCredential),
			status:   http.StatusUnauthorized,
			response: `{"message":"invalid credentials"}`,
		},
		{
			name: "error to sign token ",
			setupHandler: func() HTTPHandler {
				authenticator := &mocks.Authenticator{}
				signer := &mocks.Signer{}

				authenticator.On("Login", mock.Anything, validCredential.Username, validCredential.Password).Return(entities.ValidUser(), nil)
				signer.On("SignUser", mock.Anything, entities.ValidUser()).Return("", errors.New("error to generate token"))

				return NewHTTPHandler(authenticator, signer)
			},
			body:     serialize(&validCredential),
			status:   http.StatusUnauthorized,
			response: `{"message":"error to generate token"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()

			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			if tt.setupHandler == nil {
				t.Error("setup handler is nil")
			}

			view := tt.setupHandler()

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("POST", "/login", strings.NewReader(tt.body))
			req.Header.Add("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_httpHandler_register(t *testing.T) {
	type fields struct {
		auth service.Authenticator
		sign service.Signer
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &httpHandler{
				auth: tt.fields.auth,
				sign: tt.fields.sign,
			}
			handler.register(tt.args.c)
		})
	}
}

func serialize(entity *credential) string {
	raw, err := json.Marshal(entity)
	if err != nil {
		panic("error to marshall")
	}
	return string(raw)
}
