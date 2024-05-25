package util

import (
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization         Type = "AUTHORIZATION"
	BadRequest            Type = "BAD_REQUEST"
	Conflict              Type = "CONFLICT"
	Internal              Type = "INTERNAL"
	NotFound              Type = "NOT_FOUND"
	RequestEntityTooLarge Type = "REQUEST_ENTITY_TOO_LARGE"
	UnsupportedMediaType  Type = "UNSUPPORTED_MEDIA_TYPE"
	ServiceUnavailable    Type = "SERVICE_UNAVAILABLE"
)

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) HttpStatus() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case RequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case ServiceUnavailable:
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
		Type:    BadRequest,
		Message: fmt.Sprintf("BadRequest. Reason: %s", reason),
	}
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewConflict(name, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("Resource %v with value %v already exists", name, value),
	}
}

func NewInternal(reason string) *Error {
	return &Error{
		Type:    Internal,
		Message: reason,
	}
}

func NewNotFound(name, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("Resource %v with value %v already exists", name, value),
	}
}

func NewRequestEntityTooLarge(maxBodySize, cotentLength int64) *Error {
	return &Error{
		Type:    RequestEntityTooLarge,
		Message: fmt.Sprintf("Max payload: %v, Actual payload: %v", maxBodySize, cotentLength),
	}
}

func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: "Service unavailable or timed out",
	}
}
