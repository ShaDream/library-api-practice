package error

import "net/http"

type StatusError struct {
	error      error
	StatusCode int
}

func (s StatusError) Error() string {
	return s.error.Error()
}

func NewInternalServerError(err error) StatusError {
	return StatusError{
		error:      err,
		StatusCode: http.StatusInternalServerError,
	}
}

func NewNotFoundError(err error) StatusError {
	return StatusError{
		error:      err,
		StatusCode: http.StatusNotFound,
	}
}

func NewUnprocessableEntityError(err error) StatusError {
	return StatusError{
		error:      err,
		StatusCode: http.StatusUnprocessableEntity,
	}
}
