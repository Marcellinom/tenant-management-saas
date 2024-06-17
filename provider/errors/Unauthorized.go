package errors

import "net/http"

type UnauthorizedError struct {
	DefaultError
}

func (e UnauthorizedError) Error() string {
	return e.Message()
}

func Unauthorized(code int, message string) UnauthorizedError {
	e := UnauthorizedError{}
	e.code = code
	e.status = http.StatusUnauthorized
	e.message = message
	return e
}
