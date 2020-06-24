package contract

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Error ...
type Error struct {
	Code     int      `json:"code"`
	Messages []string `json:"messages"`
}

// NewError  ...
func NewError(code int, messages ...string) *Error {
	return &Error{Code: code, Messages: messages}
}

// BusinessError  ...
func BusinessError(messages ...string) *Error {
	return &Error{Code: http.StatusConflict, Messages: messages}
}

// BusinessErrorf  ...
func BusinessErrorf(message string, args ...interface{}) *Error {
	message = fmt.Sprintf(message, args...)
	return BusinessError(message)
}

// FromValidationError ...
func FromValidationError(e error) *Error {
	validationErrors, ok := e.(validator.ValidationErrors)

	if !ok {
		return NewError(http.StatusInternalServerError)
	}

	err := NewError(http.StatusUnprocessableEntity)
	message := "invalid value for field %s"

	for _, e := range validationErrors {
		err.Messages = append(err.Messages, fmt.Sprintf(message, e.Field()))
	}

	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"Code: %d - Messages: %s",
		e.Code,
		strings.Join(e.Messages, "\n"),
	)
}
