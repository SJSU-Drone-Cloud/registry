package utils

import (
	"errors"
	"net/http"
)

var (
	BadPostError   = errors.New("Can't post to another user's page!")
	InternalServer = errors.New("Internal server error")
)

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func BadPostSubmission(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Cannot post to another user's page!"))
}
