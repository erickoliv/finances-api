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

// CategoryRoutes add category related urls inside a gin.router/engine context
func CategoryRoutes(r *gin.RouterGroup) {
	r.POST("/categories", CreateCategory)
	r.GET("/categories/:uuid", GetCategory)
	r.PUT("/categories/:uuid", UpdateCategory)
	r.DELETE("/categories/:uuid", DeleteCategory)
	r.GET("/categories", GetCategories)
}

// GetCategories return all categories
func GetCategories(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	// println("current user", user)

	categories := []model.Category{}
	query := ExtractFilters(c.Request.URL.Query())
	query.Filters["owner = ?"] = user

	base := query.Build(db.Preloads(&categories)).Find(&categories)
	if base.Error == nil {
		response := PaginatedMessage{
			Total: query.Total,
			Page:  query.Page,
			Pages: query.Pages,
			Data:  &categories,
			Limit: query.Limit,
			Count: len(categories),
		}
		c.JSON(http.StatusOK, &response)
	} else {
		c.JSON(http.StatusInternalServerError, &base.Error)
	}
}

// CreateCategory can be called to create a new instance of Category on database, using props provided via http request
func CreateCategory(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	category := model.Category{}

	c.Bind(&category)
	category.Owner = user

	if err := db.Save(&category).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorMessage{Message: err.Error()})
	} else {
		c.JSON(http.StatusCreated, &category)
	}
}

// GetCategory can be called to get a specific category from the database. The uuid used to query is part of the url
func GetCategory(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)
	category := model.Category{}

	uuid := c.Param("uuid")

	db.Where("uuid = ? AND owner = ?", uuid, user).First(&category)

	if category.IsNew() {
		c.JSON(http.StatusNotFound, common.ErrorMessage{"category not found"})
	} else {
		c.JSON(http.StatusOK, &category)
	}
}

// UpdateCategory can be called to update a specific category. The uuid used to query is part of the url
func UpdateCategory(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	current := model.Category{}
	new := model.Category{}

	// TODO: create validate function to be used for all category related validations
	if err := c.Bind(&new); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorMessage{Message: err.Error()})
		return
	}

	uuid := c.Param("uuid")
	db.Where("uuid = ? AND owner = ?", uuid, user).First(&current)

	if current.IsNew() {
		c.JSON(http.StatusNotFound, common.ErrorMessage{"category not found"})
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

// DeleteCategory can be used to logical delete a row category from the database.
func DeleteCategory(c *gin.Context) {
	db := c.MustGet(common.DB).(*gorm.DB)
	user := c.MustGet(common.LoggedUser).(uuid.UUID)

	uuid := c.Param("uuid")
	entity := model.Category{}

	affected := db.Where("uuid = ? AND owner = ?", uuid, user).Delete(&entity).RowsAffected

	if affected > 0 {
		c.Status(http.StatusNoContent)
	} else {
		msg := fmt.Sprintf("%s - ocurrencies: %d", uuid, affected)
		c.JSON(http.StatusNotFound, common.ErrorMessage{msg})
	}
}
