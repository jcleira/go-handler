package handler

import (
	"log"
	"net/http"
)

type CustomHandler func(http.ResponseWriter, *http.Request) *HTTPError

// HTTPError struct implements the Error interface to make easier return HTTP
// errors. Using this HTTPError we are able to return the error plus the HTTP
// Status Code that we want the middleware to return on the request.
type HTTPError struct {
	Err        error `json:"error"`
	StatusCode int   `json:"status_code"`
}

// Error interface implementation.
func (ce *HTTPError) Error() string {
	return ce.Err.Error()
}

// ServerHTTP is our custom implementation needed to satisfy the Handler
// interface.
func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := ch(w, req); err != nil {
		log.Println(err)

		http.Error(w, err.Error(), err.StatusCode)
	}
}
