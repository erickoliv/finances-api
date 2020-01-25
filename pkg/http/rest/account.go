package rest

import (
	"net/http"

	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountView interface {
	Router(*gin.RouterGroup)
}

type handler struct {
	repo domain.AccountRepository
}

func MakeAccountView(repo domain.AccountRepository) AccountView {
	return handler{
		repo: repo,
	}
}

var (
	accountNotFound = domain.ErrorMessage{Message: "account not found"}
)

func (view handler) Router(group *gin.RouterGroup) {
	group.GET("/accounts", view.GetAccounts)
	group.GET("/accounts/:uuid", view.GetAccount)
	group.POST("/accounts", view.CreateAccount)
	group.PUT("/accounts/:uuid", view.UpdateAccount)
	group.DELETE("/accounts/:uuid", view.DeleteAccount)
}

// GetAccounts return all accounts
func (view handler) GetAccounts(c *gin.Context) {
	user, err := extractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	query := ExtractFilters(c.Request.URL.Query())
	query.Filters["owner = ?"] = user

	result, err := view.repo.Filter(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	response := PaginatedMessage{
		Total: query.Total,
		Page:  query.Page,
		Pages: query.Pages,
		Data:  &result,
		Limit: query.Limit,
		Count: len(result),
	}

	c.JSON(http.StatusOK, &response)
}

// CreateAccount can be called to create a new instance of Account on database, using props provided via http request
func (view handler) CreateAccount(c *gin.Context) {
	user := c.MustGet(domain.LoggedUser).(uuid.UUID)
	account := &domain.Account{}

	if err := c.ShouldBind(account); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	account.Owner = user
	if err := view.repo.Save(c, account); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount can be called to get a specific account from the database. The uuid used to query is part of the url
func (view handler) GetAccount(c *gin.Context) {
	user, err := extractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	pk, err := extractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	account, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if account.IsNew() {
		c.JSON(http.StatusNotFound, accountNotFound)
		return
	}

	c.JSON(http.StatusOK, &account)
}

// UpdateAccount can be called to update a specific account. The uuid used to query is part of the url
func (view handler) UpdateAccount(c *gin.Context) {
	new := domain.Account{}

	// TODO: create validate function to be used for all account related validations
	if err := c.Bind(&new); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorMessage{Message: err.Error()})
		return
	}
	user, err := extractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	pk, err := extractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	current, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	if current.IsNew() {
		c.JSON(http.StatusNotFound, accountNotFound)
		return
	}

	current.Name = new.Name
	current.Description = new.Description

	if err := view.repo.Save(c, current); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &current)
}

// DeleteAccount can be used to logical delete a row account from the database.
func (view handler) DeleteAccount(c *gin.Context) {
	// user := c.MustGet(domain.LoggedUser).(uuid.UUID)

	pk := uuid.MustParse(c.Param("uuid")) // fix this
	// entity := domain.Account{}

	if err := view.repo.Delete(c, pk); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, domain.ErrorMessage{err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
