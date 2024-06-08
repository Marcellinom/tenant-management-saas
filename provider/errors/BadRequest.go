package errors

import "net/http"

type BadRequestError struct {
	DefaultError
	Data map[string]any
}

func (e BadRequestError) Error() string {
	return e.message
}

func (e BadRequestError) GetData() map[string]any {
	return e.Data
}

func BadRequest(code int, message string, data ...map[string]any) BadRequestError {
	e := BadRequestError{}
	e.code = code
	e.status = http.StatusBadRequest
	e.message = message
	var error_data map[string]any
	if len(data) > 0 {
		error_data = data[0]
	}
	e.Data = error_data
	return e
}
