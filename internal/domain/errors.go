package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors for the domain layer.
var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
	ErrInternal = errors.New("internal error")
)

// ValidationFieldError represents a validation error for a specific field.
type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError holds one or more field-level validation errors.
type ValidationError struct {
	Errors []ValidationFieldError `json:"errors"`
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	msgs := make([]string, len(e.Errors))
	for i, fe := range e.Errors {
		msgs[i] = fmt.Sprintf("%s: %s", fe.Field, fe.Message)
	}
	return "validation error: " + strings.Join(msgs, "; ")
}

// IsValidationError checks if the given error is a ValidationError.
func IsValidationError(err error) (*ValidationError, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve, true
	}
	return nil, false
}
