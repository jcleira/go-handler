package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler it our Handler func enhanced to early return errors.
type Handler func(http.ResponseWriter, *http.Request) *Errors

// WithLog is the primary custom middleware, it has logging responsibility for
// request's errors, but also converts the Handler to a http.Handler.
//
// Returns a http.Handler ready to be added to the touter.
func (fn Handler) WithLog() http.Handler {
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
