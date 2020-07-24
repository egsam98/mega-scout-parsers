package errors

type TransfermarktError struct {
	Url string
}

func NewTransfermarktError(url string) *TransfermarktError {
	return &TransfermarktError{Url: url}
}

func (te *TransfermarktError) Error() string {
	return "Invalid trasfermarkt.com page: " + te.url
}
