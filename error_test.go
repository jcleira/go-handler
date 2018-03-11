package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorError(t *testing.T) {
	t.Run("When calling Error() in an Error struct", func(t *testing.T) {
		cases := []struct {
			Test   string
			Errors Errors
			Result string
		}{
			{
				Test: "with a single error",
				Errors: Errors{
					Errors: []error{
						error{
							Title: "error!",
						},
					},
				},
				Result: "error #1: error!",
			},
		}

		for _, c := range cases {
			t.Run(c.Test, func(t *testing.T) {
				assert := assert.New(t)

				errorString := c.Errors.Error()

				assert.Equal(c.Result, errorString)
			})
		}
	})
}

func TestError(t *testing.T) {
	t.Run("When creating error with Error func", func(t *testing.T) {
		cases := []struct {
			Test       string
			StatusCode int
			Args       []string
		}{
			{
				Test:       "and suceed creating an error with title, detail and source",
				StatusCode: 400,
				Args:       []string{"title", "detail", "source"},
			},
			{
				Test:       "and suceed creating an error with title, detail",
				StatusCode: 404,
				Args:       []string{"title", "detail"},
			},
			{
				Test:       "and suceed creating an error with title",
				StatusCode: 500,
				Args:       []string{"title"},
			},
		}

		for _, c := range cases {
			t.Run(c.Test, func(t *testing.T) {
				assert := assert.New(t)
				errors := Error(c.StatusCode, c.Args...)

				assert.Equal(c.StatusCode, errors.StatusCode)
				if assert.Len(errors.Errors, 1) {
					err := errors.Errors[0]

					if len(c.Args) > 0 {
						assert.Equal(c.Args[0], err.Title)
					}

					if len(c.Args) > 1 {
						assert.Equal(c.Args[1], err.Detail)
					}

					if len(c.Args) > 2 {
						assert.Equal(c.Args[2], err.Source.Pointer)
					}
				}
			})
		}
	})
}
