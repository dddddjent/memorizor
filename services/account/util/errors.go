package util

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	AuthorizationError         ErrorType = "AUTHORIZATION"
	BadRequestError            ErrorType = "BAD_REQUEST"
	ConflictError              ErrorType = "CONFLICT"
	InternalError              ErrorType = "INTERNAL"
	NotFoundError              ErrorType = "NOT_FOUND"
	RequestEntityTooLargeError ErrorType = "REQUEST_ENTITY_TOO_LARGE"
	UnsupportedMediaTypeError  ErrorType = "UNSUPPORTED_MEDIA_TYPE"
	ServiceUnavailableError    ErrorType = "SERVICE_UNAVAILABLE"
)

type Error struct {
	Type    ErrorType   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) HttpStatus() int {
	switch e.Type {
	case AuthorizationError:
		return http.StatusUnauthorized
	case BadRequestError:
		return http.StatusBadRequest
	case ConflictError:
		return http.StatusConflict
	case InternalError:
		return http.StatusInternalServerError
	case NotFoundError:
		return http.StatusNotFound
	case RequestEntityTooLargeError:
		return http.StatusRequestEntityTooLarge
	case UnsupportedMediaTypeError:
		return http.StatusUnsupportedMediaType
	case ServiceUnavailableError:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func ErrorHttpStatus(e error) int {
	if err, ok := e.(*Error); ok {
		return err.HttpStatus()
	}
	return http.StatusInternalServerError
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequestError,
		Message: fmt.Sprintf("BadRequest. Reason: %s", reason),
	}
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    AuthorizationError,
		Message: reason,
	}
}

func NewConflict(name, value string) *Error {
	return &Error{
		Type:    ConflictError,
		Message: fmt.Sprintf("Resource %v with value %v already exists", name, value),
	}
}

func NewInternal(reason string) *Error {
	return &Error{
		Type:    InternalError,
		Message: reason,
	}
}

func NewNotFound(name, value string) *Error {
	return &Error{
		Type:    NotFoundError,
		Message: fmt.Sprintf("Resource %v with value %v already exists", name, value),
	}
}

func NewRequestEntityTooLarge(maxBodySize, cotentLength int64) *Error {
	return &Error{
		Type:    RequestEntityTooLargeError,
		Message: fmt.Sprintf("Max payload: %v, Actual payload: %v", maxBodySize, cotentLength),
	}
}

func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaTypeError,
		Message: reason,
	}
}

func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailableError,
		Message: "Service unavailable or timed out",
	}
}
