package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// CustomHandler it our Handler func enhanced to early return errors.
type CustomHandler func(http.ResponseWriter, *http.Request) *Errors

// WithLog is the primary custom middleware, it has logging responsibility for
// request's errors, but also converts the CustomHandler to a http.Handler.
//
// Returns a http.Handler ready to be added to the touter.
func (fn CustomHandler) WithLog() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if e := fn(w, r); e != nil {
			defer log.Println(e.Error())
			w.WriteHeader(e.StatusCode)

			if err := json.NewEncoder(w).Encode(e); err != nil {
				defer log.Println("Unable to write the JSON error response")
			}
		}
	})
}

// WithCORS is a HTTP middleware to enable CORS on the given handler.
//
// Returns a CustomHandler with CORS enabled.
func (fn CustomHandler) WithCORS() CustomHandler {
	return CustomHandler(func(w http.ResponseWriter, r *http.Request) *Errors {
		requestHeaders := r.Header.Get("Access-Control-Request-Headers")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, DELETE, GET, OPTIONS")
		w.Header().Set("Access-Control-Request-Method", "*")
		w.Header().Set("Access-Control-Allow-Headers", requestHeaders)

		if r.Method != "OPTIONS" {
			return fn(w, r)
		}

		return nil
	})
}
