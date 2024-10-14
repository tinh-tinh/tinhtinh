package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func Exception(w http.ResponseWriter, err error, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errMsg := err.Error()
	response := map[string]interface{}{}
	if strings.IndexFunc(errMsg, func(r rune) bool { return r == '\n' }) == -1 {
		response["error"] = errMsg
	} else {
		response["error"] = strings.Split(errMsg, "\n")
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}
	return nil
}

func BadRequestException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusBadRequest)
}

func UnauthorizedException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusUnauthorized)
}

func ForbiddenException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusForbidden)
}

func NotFoundException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusNotFound)
}

func NotAllowedException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusMethodNotAllowed)
}

func ConflictException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusConflict)
}

func InternalServerException(w http.ResponseWriter, err string) error {
	return Exception(w, errors.New(err), http.StatusInternalServerError)
}
