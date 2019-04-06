package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

type PaginatedMessage struct{
	Total int `json:"total"`
	Page int `json:"page"`
	Pages int `json:"pages"`
	//Next string `json:"next"`
	//Previous string `json:"previous"`
	Limit int `json:"limit"`
	Count int `json:"count"`
	Data interface{} `json:"data"`
}

type QueryParameters struct{
	Page int
	Limit int
	Sort string
	Filters  map[string]interface{}
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

func BuildEmptyResponse(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")

	if status == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(status)
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

func ErrorResponse(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		BuildJsonResponse(ErrorMessage{Message: err.Error()}, w, http.StatusInternalServerError)
	}
}

func ValidationResponse(err error, w http.ResponseWriter) {
	BuildJsonResponse(ErrorMessage{Message: err.Error()}, w, http.StatusInternalServerError)
}


func NotFoundResponse(w http.ResponseWriter) {
	BuildJsonResponse(ErrorMessage{Message: "resource not found"}, w, http.StatusNotFound)
}

func CreatedResponse(i interface{}, w http.ResponseWriter) {
	BuildJsonResponse(i, w, http.StatusCreated)
}

func DeletedResponse(w http.ResponseWriter) {
	BuildEmptyResponse(w, http.StatusNoContent)
}

func SuccessResponse(i interface{}, w http.ResponseWriter) {
	BuildJsonResponse(i, w, http.StatusOK)
}

func PaginatedResponse(i interface{}, w http.ResponseWriter) {
	BuildJsonResponse(i, w, http.StatusOK)
}
