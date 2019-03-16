package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"olivsoft/model"
)

// GetTags return all tags
func GetTags(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := [] model.Tag{}

		app.Find(&tags)

		SuccessResponse(&tags,w)
	}
}

func CreateTag(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := model.Tag{}

		json.NewDecoder(r.Body).Decode(&tag)
		defer r.Body.Close()

		if err := app.Save(&tag).Error; err != nil {
			log.Fatal(err)
			ErrorResponse(err)
		}else {
			CreatedResponse(&tag,w)
		}
	}
}

func GetTag(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := model.Tag{}
		vars := mux.Vars(r)
		uuid := vars["uuid"]

		app.Where("uuid = ?", uuid).First(&tag)

		SuccessResponse(&tag,w)
	}
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {

}

func DeleteTag(w http.ResponseWriter, r *http.Request) {

}

// func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     w.WriteHeader(http.StatusOK)
//     fmt.Fprintf(w, "Category: %v\n", vars["category"])
// }
