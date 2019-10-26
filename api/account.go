package api

import (
	"fmt"
	"net/http"

	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/google/uuid"

	"github.com/ericktm/olivsoft-golang-api/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AccountRoutes add account related urls inside a gin.router/engine context
func AccountRoutes(r *gin.RouterGroup) {
	r.POST("/accounts", CreateAccount)
	r.GET("/accounts/:uuid", GetAccount)
	r.PUT("/accounts/:uuid", UpdateAccount)
	r.DELETE("/accounts/:uuid", DeleteAccount)
	r.GET("/accounts", GetAccounts)
}

// GetAccounts return all accounts
func GetAccounts(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	// println("current user", user)

	accounts := []model.Account{}
	query := ExtractFilters(c.Request.URL.Query())
	query.Filters["owner = ?"] = user

	base := query.Build(db.Preloads(&accounts)).Find(&accounts)
	if base.Error == nil {
		response := PaginatedMessage{
			Total: query.Total,
			Page:  query.Page,
			Pages: query.Pages,
			Data:  &accounts,
			Limit: query.Limit,
			Count: len(accounts),
		}
		c.JSON(http.StatusOK, &response)
	} else {
		c.JSON(http.StatusInternalServerError, &base.Error)
	}
}

// CreateAccount can be called to create a new instance of Account on database, using props provided via http request
func CreateAccount(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	account := model.Account{}

	c.Bind(&account)
	account.Owner = user

	if err := db.Save(&account).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
	} else {
		c.JSON(http.StatusCreated, &account)
	}
}

// GetAccount can be called to get a specific account from the database. The uuid used to query is part of the url
func GetAccount(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)
	account := model.Account{}

	uuid := c.Param("uuid")

	db.Where("uuid = ? AND owner = ?", uuid, user).First(&account)

	if account.IsNew() {
		c.JSON(http.StatusNotFound, common.ErrorMessage{"account not found"})
	} else {
		c.JSON(http.StatusOK, &account)
	}
}

// UpdateAccount can be called to update a specific account. The uuid used to query is part of the url
func UpdateAccount(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	current := model.Account{}
	new := model.Account{}

	// TODO: create validate function to be used for all account related validations
	if err := c.Bind(&new); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorMessage{Message: err.Error()})
		return
	}

	uuid := c.Param("uuid")
	db.Where("uuid = ? AND owner = ?", uuid, user).First(&current)

	if current.IsNew() {
		c.JSON(http.StatusNotFound, common.ErrorMessage{"account not found"})
	} else {
		current.Name = new.Name
		current.Description = new.Description

		if err := db.Save(&current).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, &current)
		}
	}
}

// DeleteAccount can be used to logical delete a row account from the database.
func DeleteAccount(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	uuid := c.Param("uuid")
	entity := model.Account{}

	affected := db.Where("uuid = ? AND owner = ?", uuid, user).Delete(&entity).RowsAffected

	if affected > 0 {
		c.Status(http.StatusNoContent)
	} else {
		msg := fmt.Sprintf("%s - ocurrencies: %d", uuid, affected)
		c.JSON(http.StatusNotFound, common.ErrorMessage{msg})
	}
}
