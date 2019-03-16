package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func BuildJsonResponse(c interface{}, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")

	if status == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(status)
	}

	if err := json.NewEncoder(w).Encode(&c); err != nil {
		log.Fatal("parse error", err)
	}
}

func PageNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		BuildJsonResponse(":(", w, http.StatusNotFound)
	}
}
func UnauthorizedResponse(w http.ResponseWriter) {
	msg := ErrorMessage{"unauthorized"}
	BuildJsonResponse(&msg, w, http.StatusUnauthorized)
}

func ErrorResponse(i interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		BuildJsonResponse(i, w, http.StatusInternalServerError)
	}
}

func CreatedResponse(i interface{}, w http.ResponseWriter) {
	BuildJsonResponse(i, w, http.StatusCreated)
}

func SuccessResponse(i interface{}, w http.ResponseWriter) {
	BuildJsonResponse(i, w, http.StatusOK)
}
