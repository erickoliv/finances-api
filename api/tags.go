package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"net/http"
	"olivsoft/model"
	"strconv"
	"strings"
)

// GetTags return all tags
func GetTags(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := [] model.Tag{}
		total := 0

		// TODO: put pagination handler inside mux context
		limit, _ := strconv.Atoi(r.FormValue("limit"))
		if limit == 0 {
			limit = 100
		}

		page, _ := strconv.Atoi(r.FormValue("page"))
		if page == 0 {
			page = 1
		}

		sort := r.FormValue("sort")

		// TODO: Create Generic Midleware to put filters inside context
		filt := map[string]interface{}{}
		for key, _ := range r.Form {
			if strings.HasPrefix(key, "q_") {
				if strings.HasSuffix(key, "__like") {
					field := fmt.Sprintf("%s LIKE ?", key[2:len(key)-6])
					filt[field] = r.FormValue(key)
					continue
				}
				if strings.HasSuffix(key, "__eq") {
					field := fmt.Sprintf("%s = ?", key[2:len(key)-4])
					filt[field] = r.FormValue(key)
					continue
				}
				//if strings.HasSuffix(key, "__null") {
				//	field := fmt.Sprintf("%s <= ?", key[2:len(key)-5])
				//	filt[field] = r.FormValue(key)
				//	continue
				//}
				if strings.HasSuffix(key, "__gte") {
					field := fmt.Sprintf("%s >= ?", key[2:len(key)-5])
					filt[field] = r.FormValue(key)
					continue
				}
				if strings.HasSuffix(key, "__lte") {
					field := fmt.Sprintf("%s <= ?", key[2:len(key)-5])
					filt[field] = r.FormValue(key)
					continue
				}
			}
		}

		base := app.Find(&tags)

		for k, v := range filt {
			log.Println(k, v)
			base = base.Where(k, v)
		}

		base.Count(&total)
		pages := math.Ceil(float64(total) / float64(limit))

		base = base.Offset(limit * (page - 1)).Limit(limit).Order(sort).Find(&tags)

		if base.Error == nil {
			response := PaginatedMessage{
				Total: total,
				Page:  page,
				Pages: int(pages),
				Data:  &tags,
				Limit: limit,
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
