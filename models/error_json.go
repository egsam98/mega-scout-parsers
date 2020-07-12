package models

type ErrorJSON struct {
	Error string `json:"error"`
}

func NewErrorJSON(msg string) *ErrorJSON {
	return &ErrorJSON{msg}
}
