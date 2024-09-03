package core

import (
	"encoding/json"
	"errors"
	"net/http"
)

func Exception(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Map{
		"error": err.Error(),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}

func BadRequestException(w http.ResponseWriter, err string) {
	Exception(w, errors.New(err), http.StatusBadRequest)
}

func UnauthorizedException(w http.ResponseWriter, err string) {
	Exception(w, errors.New(err), http.StatusUnauthorized)
}

func ForbiddenException(w http.ResponseWriter, err string) {
	Exception(w, errors.New(err), http.StatusForbidden)
}

func NotFoundException(w http.ResponseWriter, err string) {
	Exception(w, errors.New(err), http.StatusNotFound)
}

func ConflictException(w http.ResponseWriter, err string) {
	Exception(w, errors.New(err), http.StatusConflict)
}

func InternalServerException(w http.ResponseWriter, err string) {
	Exception(w, errors.New(err), http.StatusInternalServerError)
}
