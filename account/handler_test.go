package account

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/repository"
	"github.com/erickoliv/finances-api/test/entities"
	"github.com/erickoliv/finances-api/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeAccountView(t *testing.T) {
	mocked := &mocks.AccountService{}
	tests := []struct {
		name string
		repo repository.AccountService
		want AccountView
	}{
		{
			name: "create account view",
			repo: mocked,
			want: handler{
				repo: mocked,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MakeAccountView(tt.repo))
		})
	}
}

func Test_handler_Router(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		group *gin.RouterGroup
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
			view := handler{
				repo: tt.fields.repo,
			}
			view.Router(tt.args.group)
		})
	}
}

func Test_handler_GetAccounts(t *testing.T) {
	randomUser, _ := uuid.NewRandom()

	tests := []struct {
		name         string
		setupRepo    func() repository.AccountService
		setupContext func(c *gin.Context)
		status       int
		response     string
	}{
		{
			name: "Should return a default paginated response",
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser)
				c.Next()
			},
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Query", mock.Anything, &rest.Query{
					Page:  1,
					Limit: 100,
					Filters: map[string]interface{}{
						"owner = ?": randomUser,
					},
				}).Return(entities.ValidAcccounts(), nil)
				return repo
			},
			status:   http.StatusOK,
			response: `"page":1,"count":3`,
		},
		{
			name: "Should return a error to query",
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser)
				c.Next()
			},
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Query", mock.Anything, &rest.Query{
					Page:  1,
					Limit: 100,
					Filters: map[string]interface{}{
						"owner = ?": randomUser,
					},
				}).Return(nil, errors.New("error to query"))
				return repo
			},
			status:   http.StatusInternalServerError,
			response: `{"message":"error to query"}`,
		},
		{
			name: "Should fail if receives a context without user",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			status:   http.StatusBadRequest,
			response: `{"message":"user not present in context"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// gin.SetMode(gin.TestMode)

			router := gin.Default()

			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			view := handler{
				repo: tt.setupRepo(),
			}

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("GET", "/accounts", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_CreateAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
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
			view := handler{
				repo: tt.fields.repo,
			}
			view.CreateAccount(tt.args.c)
		})
	}
}

func Test_handler_GetAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
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
			view := handler{
				repo: tt.fields.repo,
			}
			view.GetAccount(tt.args.c)
		})
	}
}

func Test_handler_UpdateAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
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
			view := handler{
				repo: tt.fields.repo,
			}
			view.UpdateAccount(tt.args.c)
		})
	}
}

func Test_handler_DeleteAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
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
			view := handler{
				repo: tt.fields.repo,
			}
			view.DeleteAccount(tt.args.c)
		})
	}
}
