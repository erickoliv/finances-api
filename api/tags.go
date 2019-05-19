package api

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func ExtractFilters(f gin.Params) QueryParameters {

	println("parameters", f)
	// TODO: put pagination handler inside mux context
	limit, _ := strconv.Atoi(f.ByName("limit"))
	if limit == 0 {
		limit = 100
	}

	page, _ := strconv.Atoi(f.ByName("page"))
	if page == 0 {
		page = 1
	}

	sort := f.ByName("sort")

	// TODO: Create Generic Midleware to put filters inside context
	filters := map[string]interface{}{}
	// for key := range f {
	// 	if strings.HasPrefix(key, "q_") {
	// 		if strings.HasSuffix(key, "__like") {
	// 			field := fmt.Sprintf("%s LIKE ?", key[2:len(key)-6])
	// 			filters[field] = f.Get(key)
	// 			continue
	// 		}
	// 		if strings.HasSuffix(key, "__eq") {
	// 			field := fmt.Sprintf("%s = ?", key[2:len(key)-4])
	// 			filters[field] = f.Get(key)
	// 			continue
	// 		}
	// 		if strings.HasSuffix(key, "__gte") {
	// 			field := fmt.Sprintf("%s >= ?", key[2:len(key)-5])
	// 			filters[field] = f.Get(key)
	// 			continue
	// 		}
	// 		if strings.HasSuffix(key, "__lte") {
	// 			field := fmt.Sprintf("%s <= ?", key[2:len(key)-5])
	// 			filters[field] = f.Get(key)
	// 			continue
	// 		}
	// 	}
	// }

	return QueryParameters{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Filters: filters,
	}
}

// GetTags return all tags
func GetTags(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags := []model.Tag{}
		total := 0

		queryParams := ExtractFilters(c.Params)
		base := app.Preloads(&tags)

		for k, v := range queryParams.Filters {
			log.Println(k, v)
			base = base.Where(k, v)
		}

		base.Count(&total)
		pages := math.Ceil(float64(total) / float64(queryParams.Limit))

		base = base.Offset(queryParams.Limit * (queryParams.Page - 1)).Limit(queryParams.Limit).Order(queryParams.Sort).Find(&tags)

		if base.Error == nil {
			response := PaginatedMessage{
				Total: total,
				Page:  queryParams.Page,
				Pages: int(pages),
				Data:  &tags,
				Limit: queryParams.Limit,
				Count: len(tags),
			}
			c.JSON(http.StatusOK, &response)
		} else {
			c.JSON(http.StatusInternalServerError, &base.Error)
		}
	}
}

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

func UpdateTag(app *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		current := model.Tag{}
		new := model.Tag{}

		c.Bind(&new)

		uuid := c.Param("uuid")
		app.Where("uuid = ?", uuid).First(&current)

		if current.IsNew() {
			c.JSON(http.StatusNotFound, ErrorMessage{"tag not found"})
		} else {
			current.Name = new.Name
			current.Description = new.Description

			app.Save(&current)
			c.JSON(http.StatusOK, &current)
		}
	}
}

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
