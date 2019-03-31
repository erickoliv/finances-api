package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"net/http"
	"olivsoft/model"
	"strconv"
)

// GetTags return all tags
func GetTags(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := [] model.Tag{}
		total := 0

		// TODO: put pagination handler inside mux context
		limit, _ := strconv.Atoi(r.FormValue("limit"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		sort := r.FormValue("sort")


		base := app.Find(&tags)
		base.Count(&total)
		pages := math.Ceil(float64(total)/float64(limit))
		
		base = base.Offset(limit * (page-1)).Limit(limit).Order(sort).Find(&tags)

		response := PaginatedMessage{
			Total: total,
			Page: page,
			Pages: int(pages),
			Data: &tags,
			Limit: limit,
			Count: len(tags),
		}

		PaginatedResponse(&response, w)
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
