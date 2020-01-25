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

// // EntryTagRoutes add entryTag related urls inside a gin.router/engine context
// func EntryTagRoutes(r *gin.RouterGroup) {
// 	r.POST("/entry-tags", CreateEntryTag)
// 	r.GET("/entry-tags/:uuid", GetEntryTag)
// 	r.DELETE("/entry-tags/:uuid", DeleteEntryTag)
// 	r.GET("/entry-tags", GetEntryTags)
// }

// // GetEntryTags return all entry-tags
// func GetEntryTags(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	// println("current user", user)

// 	entrytags := []domain.EntryTag{}
// 	query := ExtractFilters(c.Request.URL.Query())
// 	query.Filters["owner = ?"] = user

// 	base := query.Build(db.Preloads(&entrytags)).Find(&entrytags)
// 	if base.Error == nil {
// 		response := PaginatedMessage{
// 			Total: query.Total,
// 			Page:  query.Page,
// 			Pages: query.Pages,
// 			Data:  &entrytags,
// 			Limit: query.Limit,
// 			Count: len(entrytags),
// 		}
// 		c.JSON(http.StatusOK, &response)
// 	} else {
// 		c.JSON(http.StatusInternalServerError, &base.Error)
// 	}
// }

// // CreateEntryTag can be called to create a new instance of EntryTag on database, using props provided via http request
// func CreateEntryTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	entryTag := domain.EntryTag{}

// 	if err := c.Bind(&entryTag); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrorMessage{Message: err.Error()})
// 		return
// 	}

// 	if err := db.Save(&entryTag).Error; err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrorMessage{Message: err.Error()})
// 	} else {
// 		c.JSON(http.StatusCreated, &entryTag)
// 	}
// }

// // GetEntryTag can be called to get a specific entryTag from the database. The uuid used to query is part of the url
// func GetEntryTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)
// 	entryTag := domain.EntryTag{}

// 	uuid := c.Param("uuid")

// 	db.Where("uuid = ? AND owner = ?", uuid, user).First(&entryTag)

// 	if entryTag.IsNew() {
// 		c.JSON(http.StatusNotFound, rest.ErrorMessage{"entryTag not found"})
// 	} else {
// 		c.JSON(http.StatusOK, &entryTag)
// 	}
// }

// // DeleteEntryTag can be used to logical delete a row entryTag from the database.
// func DeleteEntryTag(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	uuid := c.Param("uuid")
// 	entity := domain.EntryTag{}

// 	affected := db.Where("uuid = ? AND owner = ?", uuid, user).Delete(&entity).RowsAffected

// 	if affected > 0 {
// 		c.Status(http.StatusNoContent)
// 	} else {
// 		msg := fmt.Sprintf("%s - ocurrencies: %d", uuid, affected)
// 		c.JSON(http.StatusNotFound, rest.ErrorMessage{msg})
// 	}
// }
