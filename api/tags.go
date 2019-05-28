package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetTags return all tags
func GetTags(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags := []model.Tag{}
		query := ExtractFilters(c.Request.URL.Query())

		base := query.Build(app.Preloads(&tags)).Find(&tags)
		if base.Error == nil {
			response := PaginatedMessage{
				Total: query.Total,
				Page:  query.Page,
				Pages: query.Pages,
				Data:  &tags,
				Limit: query.Limit,
				Count: len(tags),
			}
			c.JSON(http.StatusOK, &response)
		} else {
			c.JSON(http.StatusInternalServerError, &base.Error)
		}
	}
}

// CreateTag can be called to create a new instance of Tag on database, using props provided via http request
func CreateTag(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := model.Tag{}
		c.Bind(&tag)

		if err := app.Save(&tag).Error; err != nil {
			log.Println(err)
		} else {
			c.JSON(http.StatusCreated, &tag)
		}
	}
}

// GetTag can be called to get a specific tag from the database. The uuid used to query is part of the url
func GetTag(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		tag := model.Tag{}

		uuid := c.Param("uuid")

		app.Where("uuid = ?", uuid).First(&tag)

		if tag.IsNew() {
			c.JSON(http.StatusNotFound, ErrorMessage{"tag not found"})
		} else {
			c.JSON(http.StatusOK, &tag)
		}
	}
}

// UpdateTag can be called to update a specific tag. The uuid used to query is part of the url
func UpdateTag(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		current := model.Tag{}
		new := model.Tag{}

		// TODO: create validate function to be used for all tag related validations
		if err := c.Bind(&new); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMessage{Message: err.Error()})
			return
		}

		uuid := c.Param("uuid")
		app.Where("uuid = ?", uuid).First(&current)

		if current.IsNew() {
			c.JSON(http.StatusNotFound, ErrorMessage{"tag not found"})
		} else {
			current.Name = new.Name
			current.Description = new.Description

			if err := app.Save(&current).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorMessage{Message: err.Error()})
			} else {
				c.JSON(http.StatusOK, &current)
			}
		}
	}
}

// DeleteTag can be used to logical delete a row tag from the database.
func DeleteTag(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		entity := model.Tag{}

		affected := app.Where("uuid = ?", uuid).Delete(&entity).RowsAffected

		if affected > 0 {
			c.Status(http.StatusNoContent)
		} else {
			msg := fmt.Sprintf("%s - ocurrencies: %d", uuid, affected)
			c.JSON(http.StatusNotFound, ErrorMessage{msg})
		}
	}
}
