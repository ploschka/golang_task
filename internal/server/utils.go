package server

import (
	"errors"
	"net/http"

	log "github.com/ploschka/golang_task/internal/logger"
)

var ErrInvalidPassSerie error = errors.New("INVALID PassSerie")
var ErrInvalidPassNum error = errors.New("INVALID PassNum")
var ErrInvalidPage error = errors.New("INVALID page")
var ErrInvalidLen error = errors.New("INVALID len")
var ErrJson error = errors.New("JSON ERROR")
var ErrDatabase error = errors.New("DATABASE ERROR")
var ErrBodyRead error = errors.New("REQUEST BODY READ ERROR")

func InternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Error(err)
	log.Info(http.StatusInternalServerError, "Internal server error")
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Error(err)
	log.Info(http.StatusBadRequest, "Bad request")
}

func OK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	log.Info(http.StatusOK, "Ok")
}
