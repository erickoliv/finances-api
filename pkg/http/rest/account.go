package rest

import (
	"net/http"

	"github.com/ericktm/olivsoft-golang-api/pkg/domain"
	"github.com/google/uuid"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AccountView struct {
	repo domain.AccountRepository
}

func MakeAccountView(repo domain.AccountRepository) AccountView {
	return AccountView{
		repo: repo,
	}
}

// GetAccounts return all accounts
func (view AccountView) GetAccounts(c *gin.Context) {
	user := c.MustGet(common.LoggedUser).(uuid.UUID)
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
func (view AccountView) CreateAccount(c *gin.Context) {
	user := c.MustGet(common.LoggedUser).(uuid.UUID)
	account := domain.Account{}
	c.Bind(&account)
	account.Owner = user

	if err := view.repo.Save(c, account); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &account)
}

// GetAccount can be called to get a specific account from the database. The uuid used to query is part of the url
func (view AccountView) GetAccount(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)
	account := domain.Account{}

	uuid := c.Param("uuid")

	db.Where("uuid = ? AND owner = ?", uuid, user).First(&account)

	if account.IsNew() {
		c.JSON(http.StatusNotFound, common.ErrorMessage{"account not found"})
	} else {
		c.JSON(http.StatusOK, &account)
	}
}

// UpdateAccount can be called to update a specific account. The uuid used to query is part of the url
func (view AccountView) UpdateAccount(c *gin.Context) {
	new := domain.Account{}

	// TODO: create validate function to be used for all account related validations
	if err := c.Bind(&new); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorMessage{Message: err.Error()})
		return
	}
	// _ := c.MustGet(common.LoggedUser).(uuid.UUID)

	pk, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
		return
	}

	current, err := view.repo.Get(c, pk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
		return
	}

	if current.IsNew() {
		c.JSON(http.StatusNotFound, common.ErrorMessage{"account not found"})
		return
	}

	current.Name = new.Name
	current.Description = new.Description

	if err := view.repo.Save(c, current); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &current)
}

// DeleteAccount can be used to logical delete a row account from the database.
func (view AccountView) DeleteAccount(c *gin.Context) {
	// user := c.MustGet(common.LoggedUser).(uuid.UUID)

	pk := uuid.MustParse(c.Param("uuid")) // fix this
	// entity := domain.Account{}

	if err := view.repo.Delete(c, pk); err != nil {
		c.JSON(http.StatusNotFound, common.ErrorMessage{err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
