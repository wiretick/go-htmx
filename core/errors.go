package core

import "net/http"

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	// Need this Error() to make APIError compatible with the error interface
	return e.Msg
}

func InvalidRequestError(err string) APIError {
	return APIError{
		Status: http.StatusUnprocessableEntity,
		Msg:    err,
	}
}

func NotFoundError(err string) APIError {
	return APIError{
		Status: http.StatusNotFound,
		Msg:    err,
	}
}
