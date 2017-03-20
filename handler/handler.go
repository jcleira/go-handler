package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CustomHandler it our Handler func enhanced to early return errors.
type CustomHandler func(http.ResponseWriter, *http.Request) *HTTPError

// HTTPError struct implements the Error interface to enable early return HTTP
// errors.
//
// Using an interface{} might be a little weird, this interface will receive
// error or map[string]string, the interface{} was created to make the
// HTTPError simpler to use.
//
// Examples:
//
// - Error var.
// &handler.HTTPError{err, http.StatusInternalServerError}
//
// - Error array
// var reqErrors = make(map[string]string)
//
// reqErrors["_error"] = "Unable to do that"
// reqErrors["code"] = "Code can not be null"
//
// return &handler.HTTPError{errors, http.StatusBadRequest}
//
// Used this approach to create this kind of errors:
// {
//   "errors": {
//     "_error": "Non valid Snapshot sent",
//     "cups": "Must be greater than or equal to 1",
//     "wins": "Must be greater than or equal to 1"
//   },
//   "status_code": 400
// }
//
type HTTPError struct {
	Err        interface{} `json:"error"`
	StatusCode int         `json:"status_code"`
}

// Error interface implementation.
func (ce *HTTPError) Error() string {
	errString := ""

	switch v := ce.Err.(type) {
	case error:
		errString = v.Error()
	case map[string]string:
		for key, err := range v {
			errString = fmt.Sprint(errString, " ", key, " ", err)
		}
	}

	return errString
}

// MarshalJSON implementation needed to change the default marshalling
// error's behaviour.
//
// Errors won't get marshalled, this behaviour happens because in order to do
// that the error package would have to to import the encoding/json one, by
// doing so a import cycle will happen as the encoding/json package
// use errors.
//
// You can find more information on the following issue:
// https://github.com/golang/go/issues/10748
//
// If the maps[string]string error format is used then it just throw it to the
// json.Marshl func.
func (ce *HTTPError) MarshalJSON() ([]byte, error) {
	errors := make(map[string]string)

	switch v := ce.Err.(type) {
	case error:
		errors["_error"] = v.Error()
	case map[string]string:
		errors = v
	}

	return json.Marshal(&struct {
		Errors     map[string]string `json:"errors"`
		StatusCode int               `json:"status_code"`
	}{
		Errors:     errors,
		StatusCode: ce.StatusCode,
	})
}

// ServerHTTP is our custom implementation needed to satisfy the Handler
// interface.
func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	if req.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := ch(w, req); err != nil {
		log.Println(err)

		jsonErr, errMarshal := json.Marshal(err)
		if errMarshal != nil {
			http.Error(w, "Unable to return error info, just the status code", err.StatusCode)
		}

		w.WriteHeader(err.StatusCode)
		w.Write(jsonErr)
	}
}
