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

// // EntryRoutes add entry related urls inside a gin.router/engine context
// func EntryRoutes(r *gin.RouterGroup) {
// 	r.POST("/entries", CreateEntry)
// 	r.GET("/entries/:uuid", GetEntry)
// 	r.PUT("/entries/:uuid", UpdateEntry)
// 	r.DELETE("/entries/:uuid", DeleteEntry)
// 	r.GET("/entries", GetEntries)
// }

// // GetEntries return all entries
// func GetEntries(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	// println("current user", user)

// 	entries := []domain.Entry{}
// 	query := ExtractFilters(c.Request.URL.Query())
// 	query.Filters["owner = ?"] = user

// 	base := query.Build(db.Preloads(&entries)).Find(&entries)
// 	if base.Error == nil {
// 		response := PaginatedMessage{
// 			Total: query.Total,
// 			Page:  query.Page,
// 			Pages: query.Pages,
// 			Data:  &entries,
// 			Limit: query.Limit,
// 			Count: len(entries),
// 		}
// 		c.JSON(http.StatusOK, &response)
// 	} else {
// 		c.JSON(http.StatusInternalServerError, &base.Error)
// 	}
// }

// // CreateEntry can be called to create a new instance of Entry on database, using props provided via http request
// func CreateEntry(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	entry := domain.Entry{}

// 	if err := c.Bind(&entry); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrorMessage{Message: err.Error()})
// 		return
// 	}
// 	entry.Owner = user

// 	if err := db.Save(&entry).Error; err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrorMessage{Message: err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, &entry)
// }

// // GetEntry can be called to get a specific entry from the database. The uuid used to query is part of the url
// func GetEntry(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)
// 	entry := domain.Entry{}

// 	uuid := c.Param("uuid")

// 	db.Where("uuid = ? AND owner = ?", uuid, user).First(&entry)

// 	if entry.IsNew() {
// 		c.JSON(http.StatusNotFound, rest.ErrorMessage{"entry not found"})
// 	} else {
// 		c.JSON(http.StatusOK, &entry)
// 	}
// }

// // UpdateEntry can be called to update a specific entry. The uuid used to query is part of the url
// func UpdateEntry(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	current := domain.Entry{}
// 	new := domain.Entry{}

// 	// TODO: create validate function to be used for all entry related validations
// 	if err := c.Bind(&new); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, rest.ErrorMessage{Message: err.Error()})
// 		return
// 	}

// 	uuid := c.Param("uuid")
// 	db.Where("uuid = ? AND owner = ?", uuid, user).First(&current)

// 	if current.IsNew() {
// 		c.JSON(http.StatusNotFound, rest.ErrorMessage{"entry not found"})
// 	} else {
// 		current.Name = new.Name
// 		current.Description = new.Description

// 		if err := db.Save(&current).Error; err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrorMessage{Message: err.Error()})
// 		} else {
// 			c.JSON(http.StatusOK, &current)
// 		}
// 	}
// }

// // DeleteEntry can be used to logical delete a row entry from the database.
// func DeleteEntry(c *gin.Context) {
// 	db := c.MustGet(domain.DB).(*gorm.DB)
// 	user := c.MustGet(domain.LoggedUser).(uuid.UUID)

// 	uuid := c.Param("uuid")
// 	entity := domain.Entry{}

// 	affected := db.Where("uuid = ? AND owner = ?", uuid, user).Delete(&entity).RowsAffected

// 	if affected > 0 {
// 		c.Status(http.StatusNoContent)
// 	} else {
// 		msg := fmt.Sprintf("%s - ocurrencies: %d", uuid, affected)
// 		c.JSON(http.StatusNotFound, rest.ErrorMessage{msg})
// 	}
// }
