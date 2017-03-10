package handler

import "net/http"

type CustomHandler func(http.ResponseWriter, *http.Request) error

// ServerHTTP is our custom implementation needed to satisfy the Handler
// interface.
func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := ch(w, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
