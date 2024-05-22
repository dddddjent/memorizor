package main

import (
	"net/http"
)

type Type string

const (
	Authorization         Type = "AUTHORIZATION"
	BadRequest            Type = "BADREQUEST"
	Conflict              Type = "CONFLICT"
	Internal              Type = "INTERNAL"
	NotFound              Type = "NOTFOUND"
	RequestEntityTooLarge Type = "REQUESTENTITYTOOLARGE"
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
