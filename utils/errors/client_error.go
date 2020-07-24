package errors

type ClientError struct {
	Code int
	text string
}

func NewClientError(code int, text string) *ClientError {
	return &ClientError{code, text}
}

func (ce *ClientError) Error() string {
	return ce.text
}
