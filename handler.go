package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/jsonapi"
)

// Handler is our custom Handler func enhanced to do early return for errors.
type Handler func(http.ResponseWriter, *http.Request) *jsonapi.ErrorsPayload

// ServeHTTP func that will make our custom Handler compliant with the
// http.Handler interface.
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if errors := fn(w, r); errors != nil {
		for i := 0; i < len(errors.Errors); i++ {
			log.Println(errors.Errors[i].Error())
		}

		if err := json.NewEncoder(w).Encode(errors); err != nil {
			defer log.Println("Unable to write the JSON error response")
		}
	}
}
