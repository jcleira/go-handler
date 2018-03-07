package handler

import (
	"fmt"
	"strconv"
)

// Errors is the definition created to return API JSON errors, it has been
// created using the JSON API specs, check jsonapi.org/examples/#error-objects.
//
// Also Errors struct implements the Error interface to properly log errors on
// the service log system.
//
// Example of an API error response:
//
// HTTP/1.1 400 Bad Request
// {
//   "errors": [
//     {
//       "status": "403",
//       "source": { "pointer": "/data/attributes/secret-powers" },
//       "detail": "Editing secret powers is not authorized on Sundays."
//     },
//     {
//       "status": "422",
//       "source": { "pointer": "/data/attributes/volume" },
//       "detail": "Volume does not, in fact, go to 11."
//     },
//     {
//       "status": "500",
//       "source": { "pointer": "/data/attributes/reputation" },
//       "title": "The backend responded with an error",
//       "detail": "Reputation service not responding after three requests."
//     }
//   ]
// }
//
type Errors struct {
	Errors     []error `json:"errors"`
	StatusCode int
}

// error is the definition of a single API error to return.
type error struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Source Source `json:"source"`
}

// Source is an Error attribute that points to which part caused the error.
type Source struct {
	Pointer string `json:"pointer"`
}

// Error is the Error interface implementation.
//
// Returns a unified error string.
func (e *Errors) Error() string {
	var result string
	for i, err := range e.Errors {
		result = fmt.Sprintf("%serror #%d: %s", result, i+1, err.Title)
	}

	return result
}

// Error is a package Helper func to simplify the error return of a single
// error. The objective for Error is to allow this kind of early return:
//
// Examples:
//
// handler.Error(http.StatusBadRequest, "bad body sent!")
// handler.Error(http.StatusBadRequest, "bad body sent!", err)
// handler.Error(http.StatusBadRequest, "bad body sent!", err, "/data/attributes/volume")
//
// statusCode: The request error status code to return.
// args: A ordered variadic string list of error attributes.
//
// Returns an Errors struct ready to be returned.
func Error(statusCode int, args ...string) *Errors {
	err := error{
		Status: strconv.Itoa(statusCode),
	}

	if len(args) > 0 {
		err.Title = args[0]
	}
	if len(args) > 1 {
		err.Detail = args[1]
	}
	if len(args) > 2 {
		err.Source = Source{Pointer: args[2]}
	}

	return &Errors{
		StatusCode: statusCode,
		Errors:     []error{err},
	}
}
