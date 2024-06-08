package errors

// DefaultError is an error that occurs when invariant is violated
type DefaultError struct {
	code    int
	message string
	status  int
}

func (e DefaultError) Code() int {
	return e.code
}

func (e DefaultError) Message() string {
	return e.message
}

func (e DefaultError) Status() int {
	return e.status
}
