package errors

type FetchHtmlError struct {
	inner error
}

func NewFetchHtmlError(err error) *FetchHtmlError {
	return &FetchHtmlError{inner: err}
}

func (fhe *FetchHtmlError) Error() string {
	return fhe.inner.Error()
}
