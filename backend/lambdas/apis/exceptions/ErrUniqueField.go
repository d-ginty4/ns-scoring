package exceptions

type ErrBadEntity struct {
	Message string
}

func (e *ErrBadEntity) Error() string {
	return e.Message
}

func NewErrBadEntity(message string) error {
	return &ErrBadEntity{message}
}
