package entryhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erickoliv/finances-api/auth"
	"github.com/erickoliv/finances-api/entries"
	"github.com/erickoliv/finances-api/entries/mocks"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeEntryView(t *testing.T) {
	mocked := &mocks.Repository{}
	tests := []struct {
		name string
		repo entries.Repository
		want *Handler
	}{
		{
			name: "create entry http handler",
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

func Test_handler_GetEntries(t *testing.T) {
	randomUser, _ := uuid.NewRandom()

	tests := []struct {
		name         string
		setupRepo    func() entries.Repository
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
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Query", mock.Anything, &rest.Query{
					Page:  1,
					Limit: 100,
					Filters: map[string]interface{}{
						"owner = ?": randomUser,
					},
				}).Return(mocks.ValidEntries(), nil)
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
			setupRepo: func() entries.Repository {
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
			setupRepo: func() entries.Repository {
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

			req, _ := http.NewRequest("GET", "/entries", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_CreateEntry(t *testing.T) {
	randomUser, _ := uuid.NewRandom()
	tests := []struct {
		name         string
		setupRepo    func() entries.Repository
		setupContext func(c *gin.Context)
		payload      string
		status       int
		response     string
	}{
		{
			name: "return a error with invalid payload",
			setupRepo: func() entries.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  "{}",
			status:   http.StatusBadRequest,
			response: `Error:Field validation for`,
		},
		{
			name: "return a error with invalid context, without user uuid",
			setupRepo: func() entries.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Next()
			},
			payload:  mocks.ValidEntryPayload,
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error to save a valid entry",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("error to save"))
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidEntryPayload,
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "success to save a valid entry",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, randomUser.String())
				c.Next()
			},
			payload:  mocks.ValidEntryPayload,
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
			req, _ := http.NewRequest("POST", "/entries", bytes.NewReader(payload))
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

func Test_handler_GetEntry(t *testing.T) {

	loggedUser := uuid.New()
	validEntity := uuid.New()

	tests := []struct {
		name         string
		setupRepo    func() entries.Repository
		setupContext func(c *gin.Context)
		entity       string
		status       int
		response     string
	}{
		{
			name: "error to get entry if there is no user uuid in context",
			setupRepo: func() entries.Repository {
				return &mocks.Repository{}
			},
			entity:   validEntity.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url",
			setupRepo: func() entries.Repository {
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
			name: "error to get entry from repository",
			setupRepo: func() entries.Repository {
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
			name: "empty entry from repository",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(&entries.Entry{}, nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusNotFound,
			response: `{"message":"entry not found"}`,
		},
		{
			name: "success - valid entry from repository",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Get", mock.Anything, validEntity, loggedUser).Return(mocks.ValidCompleteEntry(), nil)
				return repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, loggedUser.String())
				c.Next()
			},
			entity:   validEntity.String(),
			status:   http.StatusOK,
			response: serializeEntry(mocks.ValidCompleteEntry()),
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

			req, _ := http.NewRequest("GET", "/entries/"+tt.entity, nil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)
		})
	}
}

func Test_handler_UpdateEntry(t *testing.T) {
	validEntry := mocks.ValidEntryWithoutDescription()

	tests := []struct {
		name         string
		setupRepo    func() entries.Repository
		setupContext func(c *gin.Context)
		entity       string
		payload      string
		status       int
		response     string
	}{
		{
			name: "error due to invalid payload",
			setupRepo: func() entries.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidEntryWithoutName().UUID, mocks.InvalidValidEntryWithoutName().Owner,
				).Return(validEntry, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validEntry.Owner.String())
				c.Next()
			},
			entity:   mocks.InvalidValidEntryWithoutName().UUID.String(),
			payload:  serializeEntry(mocks.InvalidValidEntryWithoutName()),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'Entry.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name: "invalid user in context",
			setupRepo: func() entries.Repository {
				return &mocks.Repository{}
			},
			entity:   validEntry.UUID.String(),
			payload:  serializeEntry(validEntry),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "invalid uuid in url parameter",
			setupRepo: func() entries.Repository {
				return &mocks.Repository{}
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validEntry.Owner.String())
				c.Next()
			},
			entity:   "invalid param",
			payload:  serializeEntry(validEntry),
			status:   http.StatusInternalServerError,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "returns error to update a entry",
			setupRepo: func() entries.Repository {
				repo := mocks.Repository{}
				repo.On("Get",
					mock.Anything, validEntry.UUID, validEntry.Owner,
				).Return(validEntry, nil)
				repo.On("Save", mock.Anything, validEntry).Return(errors.New("error to save"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validEntry.Owner.String())
				c.Next()
			},
			entity:   validEntry.UUID.String(),
			payload:  serializeEntry(validEntry),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to save"}`,
		},
		{
			name: "error to get current entry",
			setupRepo: func() entries.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidEntryWithoutName().UUID, mocks.InvalidValidEntryWithoutName().Owner,
				).Return(nil, errors.New("error to get current entry"))

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validEntry.Owner.String())
				c.Next()
			},
			entity:   validEntry.UUID.String(),
			payload:  serializeEntry(validEntry),
			status:   http.StatusInternalServerError,
			response: `{"message":"error to get current entry"}`,
		},
		{
			name: "error due to current entry not found",
			setupRepo: func() entries.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, mocks.InvalidValidEntryWithoutName().UUID, mocks.InvalidValidEntryWithoutName().Owner,
				).Return(nil, nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validEntry.Owner.String())
				c.Next()
			},
			entity:   validEntry.UUID.String(),
			payload:  serializeEntry(validEntry),
			status:   http.StatusNotFound,
			response: `{"message":"entry not found"}`,
		},
		{
			name: "success to update a entry",
			setupRepo: func() entries.Repository {
				repo := mocks.Repository{}

				repo.On("Get",
					mock.Anything, validEntry.UUID, validEntry.Owner,
				).Return(validEntry, nil)
				repo.On("Save", mock.Anything, validEntry).Return(nil)

				return &repo
			},
			setupContext: func(c *gin.Context) {
				c.Set(auth.LoggedUser, validEntry.Owner.String())
				c.Next()
			},
			entity:   validEntry.UUID.String(),
			payload:  serializeEntry(validEntry),
			status:   http.StatusOK,
			response: serializeEntry(validEntry),
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
			req, _ := http.NewRequest("PUT", "/entries/"+tt.entity, payload)
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

func Test_handler_DeleteEntry(t *testing.T) {
	loggedUser := uuid.New()
	validEntry := mocks.ValidCompleteEntry()

	tests := []struct {
		name         string
		setupRepo    func() entries.Repository
		setupContext func(c *gin.Context)
		pk           string
		status       int
		response     string
	}{
		{
			name: "successfully delete a entry",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validEntry.UUID, loggedUser).Return(nil)
				return repo
			},
			pk:     validEntry.UUID.String(),
			status: http.StatusNoContent,
		},
		{
			name: "error to get user from context",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			setupContext: func(c *gin.Context) {

			},
			pk:       validEntry.UUID.String(),
			status:   http.StatusUnauthorized,
			response: `{"message":"user not present in context"}`,
		},
		{
			name: "error due to invalid uuid url param",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				return repo
			},
			pk:       "invalid param",
			status:   http.StatusBadRequest,
			response: `{"message":"invalid uuid in context: invalid UUID length: 13"}`,
		},
		{
			name: "error to delete a entry",
			setupRepo: func() entries.Repository {
				repo := &mocks.Repository{}
				repo.On("Delete", mock.Anything, validEntry.UUID, loggedUser).Return(errors.New("error to delete entry"))
				return repo
			},
			pk:     validEntry.UUID.String(),
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

			req, _ := http.NewRequest("DELETE", "/entries/"+tt.pk, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.status, resp.Result().StatusCode)

			body, err := ioutil.ReadAll(resp.Result().Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.response)

		})
	}
}

func serializeEntry(entity *entries.Entry) string {
	raw, err := json.Marshal(entity)
	if err != nil {
		panic("error to marshall")
	}
	return string(raw)
}
