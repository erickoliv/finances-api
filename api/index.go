package api

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
)

func IndexHandler(app *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(":)")
	}
}
