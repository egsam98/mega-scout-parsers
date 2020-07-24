package errors

type TransfermarktError struct{}

func (te *TransfermarktError) Error() string {
	return "Invalid trasfermarkt.com page"
}
