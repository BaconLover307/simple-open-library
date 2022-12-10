package exception

type ConflictError struct {
	errorMsg string
}

func NewConflictError(msg string) error {
	return ConflictError{msg}
}

func (err ConflictError) Error() string {
	return err.errorMsg
}
