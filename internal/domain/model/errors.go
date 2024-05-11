package model

import (
	"errors"
	"fmt"
)

// ValidationError represents an error that occurs due to invalid data in a specific field of a struct or input form.
type ValidationError struct {
	Field   string // Field indicates the name of the struct field associated with the error.
	Message string // Message provides a description of what is wrong with the specified field.
}

// Error returns a string representation of the ValidationError, combining the field and message.
func (ve ValidationError) Error() string {
	return fmt.Sprintf("%v : %v", ve.Field, ve.Message)
}

// ValidationErrors is a slice of ValidationError, used to aggregate multiple field errors into a single error instance.
type ValidationErrors []ValidationError

// Error returns a single string that concatenates the error messages of all contained ValidationErrors.
// It is used to implement the error interface for ValidationErrors, allowing it to be returned as an error type.
func (vs ValidationErrors) Error() string {
	var vErrs error
	for _, err := range vs {
		vErrs = errors.Join(vErrs, err) // Join appends an error to an existing error, separated by a colon.
	}
	return vErrs.Error()
}
