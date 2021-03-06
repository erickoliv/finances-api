package accounthttp

import (
	"net/http"

	"github.com/erickoliv/finances-api/accounts"
	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// HTTPHandler is a interface to expose account manipulation using HTTP
type HTTPHandler interface {
	Router(*gin.RouterGroup)
}

type handler struct {
	repo accounts.Repository
}

// NewHTTPHandler receives a Account Service, returning a internal a HTTP account handler
func NewHTTPHandler(repo accounts.Repository) HTTPHandler {
	return handler{
		repo: repo,
	}
}

func (view handler) Router(group *gin.RouterGroup) {
	group.GET("/accounts", view.GetAccounts)
	group.GET("/accounts/:uuid", view.GetAccount)
	group.POST("/accounts", view.CreateAccount)
	group.PUT("/accounts/:uuid", view.UpdateAccount)
	group.DELETE("/accounts/:uuid", view.DeleteAccount)
}

// GetAccounts return all accounts
func (view handler) GetAccounts(c *gin.Context) {
	query, err := rest.ExtractFilters(c, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := view.repo.Query(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := rest.PaginatedMessage{
		Page:  query.Page,
		Data:  &result,
		Count: len(result),
	}

	c.JSON(http.StatusOK, &response)
}

// CreateAccount can be called to create a new instance of Account on database, using props provided via http request
func (view handler) CreateAccount(c *gin.Context) {
	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	account := &accounts.Account{
		Owner: user,
	}
	if err := c.ShouldBind(account); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := view.repo.Save(c, account); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount can be called to get a specific account from the database. The uuid used to query is part of the url
func (view handler) GetAccount(c *gin.Context) {
	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	pk, err := rest.ExtractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	row, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if row.UUID == uuid.Nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "account not found"})
		return
	}

	c.JSON(http.StatusOK, row)
}

// UpdateAccount can be called to update a specific account. The uuid used to query is part of the url
func (view handler) UpdateAccount(c *gin.Context) {
	newAccount := accounts.Account{}
	if err := c.Bind(&newAccount); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	pk, err := rest.ExtractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	current, err := view.repo.Get(c, pk, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if current == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "account not found"})
		return
	}

	newAccount.UUID = pk
	newAccount.Owner = user

	if err := view.repo.Save(c, &newAccount); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, current)
}

// DeleteAccount can be used to logical delete a row account from the database.
func (view handler) DeleteAccount(c *gin.Context) {
	user, err := rest.ExtractUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	pk, err := rest.ExtractUUID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := view.repo.Delete(c, pk, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
