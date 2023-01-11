package exception

type BadRequestError struct {
	errorMsg string
}

func NewBadRequestError(msg string) error {
	return BadRequestError{msg}
}

func (err BadRequestError) Error() string {
	return err.errorMsg
}
