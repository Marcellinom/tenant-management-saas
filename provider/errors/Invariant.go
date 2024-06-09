package errors

import "net/http"

type InvariantError struct {
	DefaultError
	Data []string
}

func (e InvariantError) Error() string {
	return e.Message()
}

func Invariant(code int, message string, data ...string) InvariantError {
	e := InvariantError{}
	e.code = code
	e.status = http.StatusInternalServerError
	e.message = message
	e.Data = data
	return e
}
