package exception

type NotFoundError struct {
	errorMsg string
}

func NewNotFoundError(msg string) error {
	return NotFoundError{msg}
}

func (err NotFoundError) Error() string {
	return err.errorMsg
}
