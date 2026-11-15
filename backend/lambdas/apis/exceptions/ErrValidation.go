package exceptions

type ErrValidation struct {
	Errors []string
}

func (v *ErrValidation) Error() string { return "validation failed" }
