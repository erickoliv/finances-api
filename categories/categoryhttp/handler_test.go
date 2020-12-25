package categoryhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/categories"
	"github.com/erickoliv/finances-api/categories/mocks"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeCategoryView(t *testing.T) {
	mocked := &mocks.Repository{}
	tests := []struct {
		name string
		repo categories.Repository
		want *Handler
	}{
		{
			name: "create category http handler",
			repo: mocked,
			want: &Handler{
				repo: mocked,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewHandler(tt.repo))
		})
	}
}

func Test_handler_GetCategories(t *testing.T) {
	randomUser, _ := uuid.NewRandom()

	tests := []struct {
		name         string
		setupRepo    func() categories.Repository
		setupContext func(c *gin.Context)
		status       int
		response     string
	}{
		{
			name: "Should return a default paginated response",
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Query", mock.Anything, &rest.Query{
					Page:  1,
					Limit: 100,
					Filters: map[string]interface{}{
						"owner = ?": randomUser,
					},
				}).Return(mocks.ValidCategories(), nil)
				return repo
			},
			status:   http.StatusOK,
			response: `"page":1,"count":3`,
		},
		{
			name: "Should return a error to query",
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			setupRepo: func() categories.Repository {
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
			setupRepo: func() categories.Repository {
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

			view := Handler{
				repo: tt.setupRepo(),
			}

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("GET", "/categories", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_CreateCategory(t *testing.T) {
	randomUser, _ := uuid.NewRandom()
	tests := []struct {
		name         string
		setupRepo    func() categories.Repository
		setupContext func(c *gin.Context)
		payload      string
		status       int
		response     string
	}{
		{
			name: "return a error with invalid payload",
			setupRepo: func() categories.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  "{}",
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Category.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "return a error with invalid context, without user uuid",
			setupRepo: func() categories.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  mocks.ValidCategoryPayload,
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error to save a valid category",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("error to save"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidCategoryPayload,
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "success to save a valid category",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidCategoryPayload,
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

			view := Handler{
				repo: tt.setupRepo(),
			}

			group := router.Group("")
			view.Router(group)

			payload := []byte(tt.payload)
			req, _ := http.NewRequest("POST", "/categories", bytes.NewReader(payload))
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

func Test_handler_GetCategory(t *testing.T) {

	loggedUser := uuid.New()
	validEntity := uuid.New()

	tests := []struct {
		name         string
		setupRepo    func() categories.Repository
		setupContext func(c *gin.Context)
		entity       string
		status       int
		response     string
	}{
		{
			name: "error to get category if there is no user uuid in context",
			setupRepo: func() categories.Repository {
				return &mocks.Repository{}
			},
			entity:   validEntity.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url",
			setupRepo: func() categories.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   "invalid-uuid",
			status:   http.StatusInternalServerError,
			response: `{"message":"invalid uuid in context: invalid UUID length: 12"}`,
		},
		{
			name: "error to get category from repository",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(nil, errors.New("error to get data"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get data"}`,
		},
		{
			name: "empty category from repository",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(&categories.Category{}, nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusNotFound,
			response: `{"message":"category not found"}`,
		},
		{
			name: "success - valid category from repository",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(mocks.ValidCompleteCategory(), nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusOK,
			response: serializeCategory(mocks.ValidCompleteCategory()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()
			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			view := NewHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("GET", "/categories/"+tt.entity, nil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_UpdateCategory(t *testing.T) {
	validTag := mocks.ValidCategoryWithoutDescription()

	tests := []struct {
		name         string
		setupRepo    func() categories.Repository
		setupContext func(c *gin.Context)
		entity       string
		payload      string
		status       int
		response     string
	}{
		{
			name: "error due to invalid payload",
			setupRepo: func() categories.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidCategoryWithoutName().UUID, mocks.InvalidValidCategoryWithoutName().Owner,
				).Return(validTag, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   mocks.InvalidValidCategoryWithoutName().UUID.String(),
			payload:  serializeCategory(mocks.InvalidValidCategoryWithoutName()),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Category.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "invalid user in context",
			setupRepo: func() categories.Repository {
				return &mocks.Repository{}
			},
			entity:   validTag.UUID.String(),
			payload:  serializeCategory(validTag),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url parameter",
			setupRepo: func() categories.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   "invalid param",
			payload:  serializeCategory(validTag),
			status:   http.StatusInternalServerError,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "returns error to update a category",
			setupRepo: func() categories.Repository {
				repo := mocks.Repository{}
				repo.On("Get",
					mock.Anything, validTag.UUID, validTag.Owner,
				).Return(validTag, nil)
				repo.On("Save", mock.Anything, validTag).Return(errors.New("error to save"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   validTag.UUID.String(),
			payload:  serializeCategory(validTag),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "error to get current category",
			setupRepo: func() categories.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidCategoryWithoutName().UUID, mocks.InvalidValidCategoryWithoutName().Owner,
				).Return(nil, errors.New("error to get current category"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   validTag.UUID.String(),
			payload:  serializeCategory(validTag),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get current category"}`,
		},
		{
			name: "error due to current category not found",
			setupRepo: func() categories.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidCategoryWithoutName().UUID, mocks.InvalidValidCategoryWithoutName().Owner,
				).Return(nil, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   validTag.UUID.String(),
			payload:  serializeCategory(validTag),
			status:   http.StatusNotFound,
			response: `{"message":"category not found"}`,
		},
		{
			name: "success to update a category",
			setupRepo: func() categories.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, validTag.UUID, validTag.Owner,
				).Return(validTag, nil)
				repo.On("Save", mock.Anything, validTag).Return(nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   validTag.UUID.String(),
			payload:  serializeCategory(validTag),
			status:   http.StatusOK,
			response: serializeCategory(validTag),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			}

			view := NewHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			payload := bytes.NewBufferString(tt.payload)
			req, _ := http.NewRequest("PUT", "/categories/"+tt.entity, payload)
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

func Test_handler_DeleteCategory(t *testing.T) {
	loggedUser := uuid.New()
	validTag := mocks.ValidCompleteCategory()

	tests := []struct {
		name         string
		setupRepo    func() categories.Repository
		setupContext func(c *gin.Context)
		pk           string
		status       int
		response     string
	}{
		{
			name: "successfully delete a category",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validTag.UUID, loggedUser).Return(nil)
				return repo
			},
			pk:     validTag.UUID.String(),
			status: http.StatusNoContent,
		},
		{
			name: "error to get user from context",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			setupContext: func(c *gin.Context) {

			},
			pk:       validTag.UUID.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error due to invalid uuid url param",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			pk:       "invalid param",
			status:   http.StatusBadRequest,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "error to delete a category",
			setupRepo: func() categories.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validTag.UUID, loggedUser).Return(errors.New("error to delete category"))
				return repo
			},
			pk:     validTag.UUID.String(),
			status: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.New()

			if tt.setupContext == nil {
				tt.setupContext = func(c *gin.Context) {
					c.Set(auth.LoggedUser, loggedUser.String())
					c.Next()
				}
			}

			router.Use(tt.setupContext)
			view := NewHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("DELETE", "/categories/"+tt.pk, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)

		})
	}
}

func serializeCategory(entity *categories.Category) string {
	raw, err := json.Marshal(entity)
	if err != nil {
		panic("error to marshall")
	}
	return string(raw)
}
