package exceptions

type ErrUniqueField struct {
	Message string
}

func (e *ErrUniqueField) Error() string {
	return e.Message
}

func NewErrUniqueField(message string) error {
	return &ErrUniqueField{message}
}
