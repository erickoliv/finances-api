package account

import (
	"bytes"
	"encoding/json"
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
		want HTTPHandler
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
			assert.Equal(t, tt.want, NewHTTPHandler(tt.repo))
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
	randomUser, _ := uuid.NewRandom()
	tests := []struct {
		name         string
		setupRepo    func() repository.AccountService
		setupContext func(c *gin.Context)
		payload      string
		status       int
		response     string
	}{
		{
			name: "return a error with invalid payload",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  "{}",
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "return a error with invalid context, without user uuid",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  entities.ValidAccountPayload,
			status:   http.StatusBadRequest,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error to save a valid account",
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("error to save"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser)
				c.Next()
			},
			payload:  entities.ValidAccountPayload,
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "success to save a valid account",
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser)
				c.Next()
			},
			payload:  entities.ValidAccountPayload,
			status:   http.StatusCreated,
			response: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			view := handler{
				repo: tt.setupRepo(),
			}

			group := router.Group("")
			view.Router(group)

			payload := []byte(tt.payload)
			req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(payload))
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

func Test_handler_GetAccount(t *testing.T) {

	loggedUser := uuid.New()
	validEntity := uuid.New()

	tests := []struct {
		name         string
		setupRepo    func() repository.AccountService
		setupContext func(c *gin.Context)
		entity       string
		status       int
		response     string
	}{
		{
			name: "error to get account if there is no user uuid in context",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			entity:   validEntity.String(),
			status:   http.StatusInternalServerError,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser)
				c.Next()
			},
			entity:   "invalid-uuid",
			status:   http.StatusInternalServerError,
			response: `{"message":"uuid parameter is invalid"}`,
		},
		{
			name: "error to get account from repository",
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(nil, errors.New("error to get data"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser)
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get data"}`,
		},
		{
			name: "empty account from repository",
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(&domain.Account{}, nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser)
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusNotFound,
			response: `{"message":"account not found"}`,
		},
		{
			name: "success - valid account from repository",
			setupRepo: func() repository.AccountService {
				repo := &mocks.AccountService{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(entities.ValidCompleteAccount(), nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser)
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusOK,
			response: serializeAccount(entities.ValidCompleteAccount()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			view := NewHTTPHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("GET", "/accounts/"+tt.entity, nil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_UpdateAccount(t *testing.T) {
	validAccount := entities.ValidAccountWithoutDescription()

	tests := []struct {
		name         string
		setupRepo    func() repository.AccountService
		setupContext func(c *gin.Context)
		entity       string
		payload      string
		status       int
		response     string
	}{
		{
			name: "error due to invalid payload",
			setupRepo: func() repository.AccountService {
				repo := mocks.AccountService{}

				repo.On("Get",
					mock.Anything, entities.ValidAccountWithoutName().UUID, entities.ValidAccountWithoutName().Owner,
				).Return(validAccount, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner)
				c.Next()
			},
			entity:   entities.ValidAccountWithoutName().UUID.String(),
			payload:  serializeAccount(entities.ValidAccountWithoutName()),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "invalid user in context",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url parameter",
			setupRepo: func() repository.AccountService {
				return &mocks.AccountService{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner)
				c.Next()
			},
			entity:   "invalid param",
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"uuid parameter is invalid"}`,
		},
		{
			name: "returns error to update a account",
			setupRepo: func() repository.AccountService {
				repo := mocks.AccountService{}
				repo.On("Get",
					mock.Anything, validAccount.UUID, validAccount.Owner,
				).Return(validAccount, nil)
				repo.On("Save", mock.Anything, validAccount).Return(errors.New("error to save"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner)
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "error to get current account",
			setupRepo: func() repository.AccountService {
				repo := mocks.AccountService{}

				repo.On("Get",
					mock.Anything, entities.ValidAccountWithoutName().UUID, entities.ValidAccountWithoutName().Owner,
				).Return(nil, errors.New("error to get current account"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner)
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get current account"}`,
		},
		{
			name: "error due to current account not found",
			setupRepo: func() repository.AccountService {
				repo := mocks.AccountService{}

				repo.On("Get",
					mock.Anything, entities.ValidAccountWithoutName().UUID, entities.ValidAccountWithoutName().Owner,
				).Return(nil, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner)
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusNotFound,
			response: `{"message":"account not found"}`,
		},
		{
			name: "success to update a account",
			setupRepo: func() repository.AccountService {
				repo := mocks.AccountService{}

				repo.On("Get",
					mock.Anything, validAccount.UUID, validAccount.Owner,
				).Return(validAccount, nil)
				repo.On("Save", mock.Anything, validAccount).Return(nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner)
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusOK,
			response: serializeAccount(validAccount),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			view := NewHTTPHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			payload := bytes.NewBufferString(tt.payload)
			req, _ := http.NewRequest("PUT", "/accounts/"+tt.entity, payload)
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

func serializeAccount(entity *domain.Account) string {
	raw, err := json.Marshal(entity)
	if err != nil {
		panic("error to marshall")
	}
	return string(raw)
}
