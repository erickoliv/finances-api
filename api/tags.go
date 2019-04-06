package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"net/http"
	"net/url"
	"olivsoft/model"
	"strconv"
	"strings"
)

func ExtractFilters(f url.Values) QueryParameters {

	println("parameters", f)
	// TODO: put pagination handler inside mux context
	limit, _ := strconv.Atoi(f.Get("limit"))
	if limit == 0 {
		limit = 100
	}

	page, _ := strconv.Atoi(f.Get("page"))
	if page == 0 {
		page = 1
	}

	sort := f.Get("sort")

	log.Println("lluasdasd", f)

	// TODO: Create Generic Midleware to put filters inside context
	filters := map[string]interface{}{}
	for key := range f {
		if strings.HasPrefix(key, "q_") {
			if strings.HasSuffix(key, "__like") {
				field := fmt.Sprintf("%s LIKE ?", key[2:len(key)-6])
				filters[field] = f.Get(key)
				continue
			}
			if strings.HasSuffix(key, "__eq") {
				field := fmt.Sprintf("%s = ?", key[2:len(key)-4])
				filters[field] = f.Get(key)
				continue
			}
			if strings.HasSuffix(key, "__gte") {
				field := fmt.Sprintf("%s >= ?", key[2:len(key)-5])
				filters[field] = f.Get(key)
				continue
			}
			if strings.HasSuffix(key, "__lte") {
				field := fmt.Sprintf("%s <= ?", key[2:len(key)-5])
				filters[field] = f.Get(key)
				continue
			}
		}
	}

	return QueryParameters{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Filters: filters,
	}
}

// GetTags return all tags
func GetTags(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := []model.Tag{}
		total := 0

		if err := r.ParseForm(); err != nil {
			ValidationResponse(err, w)
			return
		}
		queryParams := ExtractFilters(r.Form)
		base := app.Preloads(&tags)

		log.Println(queryParams)

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
			PaginatedResponse(&response, w)
		} else {
			ValidationResponse(base.Error, w)
		}
	}
}

func CreateTag(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := model.Tag{}

		json.NewDecoder(r.Body).Decode(&tag)
		defer r.Body.Close()

		if err := app.Save(&tag).Error; err != nil {
			log.Println(err)
			ValidationResponse(err, w)
		} else {
			CreatedResponse(&tag, w)
		}
	}
}

func GetTag(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := model.Tag{}
		vars := mux.Vars(r)
		uuid := vars["uuid"]

		app.Where("uuid = ?", uuid).First(&tag)

		if tag.IsNew() {
			NotFoundResponse(w)
		} else {
			SuccessResponse(&tag, w)
		}
	}
}

func UpdateTag(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		current := model.Tag{}
		new := model.Tag{}

		json.NewDecoder(r.Body).Decode(&new)
		defer r.Body.Close()

		vars := mux.Vars(r)
		uuid := vars["uuid"]

		app.Where("uuid = ?", uuid).First(&current)

		if current.IsNew() {
			NotFoundResponse(w)
		} else {
			current.Name = new.Name
			current.Description = new.Description

			app.Save(&current)
			SuccessResponse(&current, w)
		}
	}
}

func DeleteTag(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]

		app.Where("uuid = ?", uuid).Delete(&model.Tag{})

		DeletedResponse(w)

	}
}
