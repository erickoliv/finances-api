package rest

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/erickoliv/finances-api/domain"
// 	"github.com/google/uuid"

// 	"github.com/erickoliv/finances-api/domain"
// 	"github.com/gin-gonic/gin"
// 	"github.com/jinzhu/gorm"
// )

// // TagRoutes add tag related urls inside a gin.router/engine context
// func TagRoutes(r *gin.RouterGroup) {
// 	r.POST("/tags", CreateTag)
// 	r.GET("/tags/:uuid", GetTag)
// 	r.PUT("/tags/:uuid", UpdateTag)
// 	r.DELETE("/tags/:uuid", DeleteTag)
// 	r.GET("/tags", GetTags)
// }

// // GetTags return all tags
// func GetTags(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	// println("current user", user)

// 	tags := []domain.Tag{}
// 	query := ExtractFilters(c.Request.URL.Query())
// 	query.Filters["owner = ?"] = user

// 	base := query.Build(db.Preloads(&tags)).Find(&tags)
// 	if base.Error == nil {
// 		response := PaginatedMessage{
// 			Total: query.Total,
// 			Page:  query.Page,
// 			Pages: query.Pages,
// 			Data:  &tags,
// 			Limit: query.Limit,
// 			Count: len(tags),
// 		}
// 		c.JSON(http.StatusOK, &response)
// 	} else {
// 		c.JSON(http.StatusInternalServerError, &base.Error)
// 	}
// }

// // CreateTag can be called to create a new instance of Tag on database, using props provided via http request
// func CreateTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	tag := domain.Tag{}

// 	c.Bind(&tag)
// 	tag.Owner = user

// 	if err := db.Save(&tag).Error; err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
// 	} else {
// 		c.JSON(http.StatusCreated, &tag)
// 	}
// }

// // GetTag can be called to get a specific tag from the database. The uuid used to query is part of the url
// func GetTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)
// 	tag := domain.Tag{}

// 	uuid := c.Param("uuid")

// 	db.Where("uuid = ? AND owner = ?", uuid, user).First(&tag)

// 	if tag.IsNew() {
// 		c.JSON(http.StatusNotFound, domain.ErrorMessage{"tag not found"})
// 	} else {
// 		c.JSON(http.StatusOK, &tag)
// 	}
// }

// // UpdateTag can be called to update a specific tag. The uuid used to query is part of the url
// func UpdateTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	current := domain.Tag{}
// 	new := domain.Tag{}

// 	// TODO: create validate function to be used for all tag related validations
// 	if err := c.Bind(&new); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorMessage{Message: err.Error()})
// 		return
// 	}

// 	uuid := c.Param("uuid")
// 	db.Where("uuid = ? AND owner = ?", uuid, user).First(&current)

// 	if current.IsNew() {
// 		c.JSON(http.StatusNotFound, domain.ErrorMessage{"tag not found"})
// 	} else {
// 		current.Name = new.Name
// 		current.Description = new.Description

// 		if err := db.Save(&current).Error; err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorMessage{Message: err.Error()})
// 		} else {
// 			c.JSON(http.StatusOK, &current)
// 		}
// 	}
// }

// // DeleteTag can be used to logical delete a row tag from the database.
// func DeleteTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	uuid := c.Param("uuid")
// 	entity := domain.Tag{}

// 	affected := db.Where("uuid = ? AND owner = ?", uuid, user).Delete(&entity).RowsAffected

// 	if affected > 0 {
// 		c.Status(http.StatusNoContent)
// 	} else {
// 		msg := fmt.Sprintf("%s - ocurrencies: %d", uuid, affected)
// 		c.JSON(http.StatusNotFound, domain.ErrorMessage{msg})
// 	}
// }
