package taghttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/erickoliv/finances-api/tags"
	"github.com/erickoliv/finances-api/tags/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeTagView(t *testing.T) {
	mocked := &mocks.Repository{}
	tests := []struct {
		name string
		repo tags.Repository
		want HTTPHandler
	}{
		{
			name: "create tag view",
			repo: mocked,
			want: &Handler{
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

func Test_handler_GetTags(t *testing.T) {
	randomUser, _ := uuid.NewRandom()

	tests := []struct {
		name         string
		setupRepo    func() tags.Repository
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
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Query", mock.Anything, &rest.Query{
					Page:  1,
					Limit: 100,
					Filters: map[string]interface{}{
						"owner = ?": randomUser,
					},
				}).Return(mocks.ValidTags(), nil)
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
			setupRepo: func() tags.Repository {
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
			setupRepo: func() tags.Repository {
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

			req, _ := http.NewRequest("GET", "/tags", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_CreateTag(t *testing.T) {
	randomUser, _ := uuid.NewRandom()
	tests := []struct {
		name         string
		setupRepo    func() tags.Repository
		setupContext func(c *gin.Context)
		payload      string
		status       int
		response     string
	}{
		{
			name: "return a error with invalid payload",
			setupRepo: func() tags.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  "{}",
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Tag.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "return a error with invalid context, without user uuid",
			setupRepo: func() tags.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  mocks.ValidTagPayload,
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error to save a valid tag",
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("error to save"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidTagPayload,
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "success to save a valid tag",
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidTagPayload,
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
			req, _ := http.NewRequest("POST", "/tags", bytes.NewReader(payload))
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

func Test_handler_GetTag(t *testing.T) {

	loggedUser := uuid.New()
	validEntity := uuid.New()

	tests := []struct {
		name         string
		setupRepo    func() tags.Repository
		setupContext func(c *gin.Context)
		entity       string
		status       int
		response     string
	}{
		{
			name: "error to get tag if there is no user uuid in context",
			setupRepo: func() tags.Repository {
				return &mocks.Repository{}
			},
			entity:   validEntity.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url",
			setupRepo: func() tags.Repository {
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
			name: "error to get tag from repository",
			setupRepo: func() tags.Repository {
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
			name: "empty tag from repository",
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(&tags.Tag{}, nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusNotFound,
			response: `{"message":"tag not found"}`,
		},
		{
			name: "success - valid tag from repository",
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(mocks.ValidCompleteTag(), nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusOK,
			response: serializeTag(mocks.ValidCompleteTag()),
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

			req, _ := http.NewRequest("GET", "/tags/"+tt.entity, nil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_UpdateTag(t *testing.T) {
	validTag := mocks.ValidTagWithoutDescription()

	tests := []struct {
		name         string
		setupRepo    func() tags.Repository
		setupContext func(c *gin.Context)
		entity       string
		payload      string
		status       int
		response     string
	}{
		{
			name: "error due to invalid payload",
			setupRepo: func() tags.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidTagWithoutName().UUID, mocks.InvalidValidTagWithoutName().Owner,
				).Return(validTag, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   mocks.InvalidValidTagWithoutName().UUID.String(),
			payload:  serializeTag(mocks.InvalidValidTagWithoutName()),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Tag.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "invalid user in context",
			setupRepo: func() tags.Repository {
				return &mocks.Repository{}
			},
			entity:   validTag.UUID.String(),
			payload:  serializeTag(validTag),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url parameter",
			setupRepo: func() tags.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   "invalid param",
			payload:  serializeTag(validTag),
			status:   http.StatusInternalServerError,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "returns error to update a tag",
			setupRepo: func() tags.Repository {
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
			payload:  serializeTag(validTag),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "error to get current tag",
			setupRepo: func() tags.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidTagWithoutName().UUID, mocks.InvalidValidTagWithoutName().Owner,
				).Return(nil, errors.New("error to get current tag"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   validTag.UUID.String(),
			payload:  serializeTag(validTag),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get current tag"}`,
		},
		{
			name: "error due to current tag not found",
			setupRepo: func() tags.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidTagWithoutName().UUID, mocks.InvalidValidTagWithoutName().Owner,
				).Return(nil, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validTag.Owner.String())
				c.Next()
			},
			entity:   validTag.UUID.String(),
			payload:  serializeTag(validTag),
			status:   http.StatusNotFound,
			response: `{"message":"tag not found"}`,
		},
		{
			name: "success to update a tag",
			setupRepo: func() tags.Repository {
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
			payload:  serializeTag(validTag),
			status:   http.StatusOK,
			response: serializeTag(validTag),
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
			req, _ := http.NewRequest("PUT", "/tags/"+tt.entity, payload)
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

func Test_handler_DeleteTag(t *testing.T) {
	loggedUser := uuid.New()
	validTag := mocks.ValidCompleteTag()

	tests := []struct {
		name         string
		setupRepo    func() tags.Repository
		setupContext func(c *gin.Context)
		pk           string
		status       int
		response     string
	}{
		{
			name: "successfully delete a tag",
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validTag.UUID, loggedUser).Return(nil)
				return repo
			},
			pk:     validTag.UUID.String(),
			status: http.StatusNoContent,
		},
		{
			name: "error to get user from context",
			setupRepo: func() tags.Repository {
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
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			pk:       "invalid param",
			status:   http.StatusBadRequest,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "error to delete a tag",
			setupRepo: func() tags.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validTag.UUID, loggedUser).Return(errors.New("error to delete tag"))
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
			view := NewHTTPHandler(tt.setupRepo())

			group := router.Group("")
			view.Router(group)

			req, _ := http.NewRequest("DELETE", "/tags/"+tt.pk, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)

		})
	}
}

func serializeTag(entity *tags.Tag) string {
	raw, err := json.Marshal(entity)
	if err != nil {
		panic("error to marshall")
	}
	return string(raw)
}
