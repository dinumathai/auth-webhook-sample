package api

import (
	"net/http"
)

// JSONWriter interface
type JSONWriter interface {
	Write(w http.ResponseWriter)
}

// JSONResponse defines a JSON response sent to client
type JSONResponse struct {
	status int
	data   []byte
}

// ErrorResponse is
type ErrorResponse struct {
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (jr *JSONResponse) Write(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(jr.status)

	w.Write(jr.data)
}
