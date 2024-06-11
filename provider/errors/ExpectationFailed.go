package errors

import "net/http"

type ExpectationFailedError struct {
	DefaultError
}

func (e ExpectationFailedError) Error() string {
	return e.message
}

func ExpectationFailed(code int, message string) ExpectationFailedError {
	e := ExpectationFailedError{}
	e.status = http.StatusExpectationFailed
	e.code = code
	e.message = message
	return e
}
