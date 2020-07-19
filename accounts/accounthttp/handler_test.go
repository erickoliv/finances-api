package accounthttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/accounts"
	"github.com/erickoliv/finances-api/accounts/mocks"
	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeAccountView(t *testing.T) {
	mocked := &mocks.Repository{}
	tests := []struct {
		name string
		repo accounts.Repository
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
		setupRepo    func() accounts.Repository
		setupContext func(c *gin.Context)
		status       int
		response     string
	}{
		{
			name: "Should return a default paginated response",
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser.String())
				c.Next()
			},
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Query", mock.Anything, &rest.Query{
					Page:  1,
					Limit: 100,
					Filters: map[string]interface{}{
						"owner = ?": randomUser,
					},
				}).Return(mocks.ValidAcccounts(), nil)
				return repo
			},
			status:   http.StatusOK,
			response: `"page":1,"count":3`,
		},
		{
			name: "Should return a error to query",
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser.String())
				c.Next()
			},
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
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
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			status:   http.StatusUnauthorized,
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
		setupRepo    func() accounts.Repository
		setupContext func(c *gin.Context)
		payload      string
		status       int
		response     string
	}{
		{
			name: "return a error with invalid context, without user uuid",
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  mocks.ValidAccountPayload,
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "return a error with invalid payload",
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  "{}",
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "error to save a valid account",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("error to save"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidAccountPayload,
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "success to save a valid account",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, randomUser.String())
				c.Next()
			},
			payload: mocks.ValidAccountPayload,
			status:  http.StatusCreated,
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
		setupRepo    func() accounts.Repository
		setupContext func(c *gin.Context)
		entity       string
		status       int
		response     string
	}{
		{
			name: "error to get account if there is no user uuid in context",
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			entity:   validEntity.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url",
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   "invalid-uuid",
			status:   http.StatusInternalServerError,
			response: `{"message":"invalid uuid in context: invalid UUID length: 12"}`,
		},
		{
			name: "error to get account from repository",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(nil, errors.New("error to get data"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get data"}`,
		},
		{
			name: "empty account from repository",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(&accounts.Account{}, nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusNotFound,
			response: `{"message":"account not found"}`,
		},
		{
			name: "success - valid account from repository",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(mocks.ValidCompleteAccount(), nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusOK,
			response: serializeAccount(mocks.ValidCompleteAccount()),
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
	validAccount := mocks.ValidAccountWithoutDescription()

	tests := []struct {
		name         string
		setupRepo    func() accounts.Repository
		setupContext func(c *gin.Context)
		entity       string
		payload      string
		status       int
		response     string
	}{
		{
			name: "error due to invalid payload",
			setupRepo: func() accounts.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.ValidAccountWithoutName().UUID, mocks.ValidAccountWithoutName().Owner,
				).Return(validAccount, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner.String())
				c.Next()
			},
			entity:   mocks.ValidAccountWithoutName().UUID.String(),
			payload:  serializeAccount(mocks.ValidAccountWithoutName()),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "invalid user in context",
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url parameter",
			setupRepo: func() accounts.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner.String())
				c.Next()
			},
			entity:   "invalid param",
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "returns error to update a account",
			setupRepo: func() accounts.Repository {
				repo := mocks.Repository{}
				repo.On("Get",
					mock.Anything, validAccount.UUID, validAccount.Owner,
				).Return(validAccount, nil)
				repo.On("Save", mock.Anything, validAccount).Return(errors.New("error to save"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner.String())
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "error to get current account",
			setupRepo: func() accounts.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.ValidAccountWithoutName().UUID, mocks.ValidAccountWithoutName().Owner,
				).Return(nil, errors.New("error to get current account"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner.String())
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get current account"}`,
		},
		{
			name: "error due to current account not found",
			setupRepo: func() accounts.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.ValidAccountWithoutName().UUID, mocks.ValidAccountWithoutName().Owner,
				).Return(nil, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner.String())
				c.Next()
			},
			entity:   validAccount.UUID.String(),
			payload:  serializeAccount(validAccount),
			status:   http.StatusNotFound,
			response: `{"message":"account not found"}`,
		},
		{
			name: "success to update a account",
			setupRepo: func() accounts.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, validAccount.UUID, validAccount.Owner,
				).Return(validAccount, nil)
				repo.On("Save", mock.Anything, validAccount).Return(nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(domain.LoggedUser, validAccount.Owner.String())
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
	loggedUser := uuid.New()
	validAccount := mocks.ValidCompleteAccount()

	tests := []struct {
		name         string
		setupRepo    func() accounts.Repository
		setupContext func(c *gin.Context)
		pk           string
		status       int
		response     string
	}{
		{
			name: "successfully delete a account",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validAccount.UUID, loggedUser).Return(nil)
				return repo
			},
			pk:     validAccount.UUID.String(),
			status: http.StatusNoContent,
		},
		{
			name: "error to get user from context",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			setupContext: func(c *gin.Context) {

			},
			pk:       validAccount.UUID.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error due to invalid uuid url param",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			pk:       "invalid param",
			status:   http.StatusBadRequest,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "error to delete a account",
			setupRepo: func() accounts.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validAccount.UUID, loggedUser).Return(errors.New("error to delete account"))
				return repo
			},
			pk:     validAccount.UUID.String(),
			status: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()

			if tt.setupContext == nil {
				tt.setupContext = func(c *gin.Context) {
					c.Set(domain.LoggedUser, loggedUser.String())
					c.Next()
				}
			}

			router.Use(tt.setupContext)
			view := NewHTTPHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("DELETE", "/accounts/"+tt.pk, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)

		})
	}
}

func serializeAccount(entity *accounts.Account) string {
	raw, err := json.Marshal(entity)
	if err != nil {
		panic("error to marshall")
	}
	return string(raw)
}
