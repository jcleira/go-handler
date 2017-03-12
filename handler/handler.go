package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CustomHandler func(http.ResponseWriter, *http.Request) *HTTPError

// HTTPError struct implements the Error interface to make easier return HTTP
// errors. Using this HTTPError we are able to return the error plus the HTTP
// Status Code that we want the middleware to return on the request.
//
// Using an interface{} might be a little weird, this interface will receive
// error or []error, the interface{} was created to make the HTTPError simpler
// to use.
//
// Examples:
//
// - Error var.
// &handler.HTTPError{err, http.StatusInternalServerError}
//
// - Error array
// for _, err := range result.Errors() {
//	errors = append(errors, fmt.Errorf("%s", err))
// }
// return &handler.HTTPError{errors, http.StatusBadRequest}
//
// Used this approach to create this kind of errors:
// {
//   "errors": [
//	 	 "Non valid Snapshot sent.",
//		 "card_3: String length must be greater than or equal to 3",
//		 "card_7: String length must be greater than or equal to 3"
//	 ],
//	 "status_code": 400
// }
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
	case []error:
		for _, err := range v {
			errString = fmt.Sprint(errString, " ", err)
		}
	}

	return errString
}

// MarshalJSON implementation needed to change the default Marshall behaviour
// with errors.
//
// Errors won't get Marshalled by themselves, this behaviour occurs because in
// order to be able to Marshal errors the error package would have to to import
// the encoding/json package, creating a import cycle as obviously the
// encoding/json uses errors.
//
// You can find more information on the following issue:
// https://github.com/golang/go/issues/10748
func (ce *HTTPError) MarshalJSON() ([]byte, error) {
	var errors []string

	switch v := ce.Err.(type) {
	case error:
		errors = append(errors, v.Error())
	case []error:
		for _, err := range v {
			errors = append(errors, err.Error())
		}
	}

	return json.Marshal(&struct {
		Errors     []string `json:"errors"`
		StatusCode int      `json:"status_code"`
	}{
		Errors:     errors,
		StatusCode: ce.StatusCode,
	})
}

// ServerHTTP is our custom implementation needed to satisfy the Handler
// interface.
func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := ch(w, req); err != nil {
		log.Println(err)

		w.WriteHeader(err.StatusCode)

		jsonErr, err := json.Marshal(err)
		if err != nil {
			log.Println("Error while Marshaling the HTTPError: ", err)
			return
		}

		w.Write(jsonErr)
	}
}
