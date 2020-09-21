package core

// ErrValidation represents a validation error.
type ErrValidation struct {
	Message string
	Err     error
}

func (e ErrValidation) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewErrValidation creates a new ErrValidation.
func NewErrValidation(message string) ErrValidation {
	return ErrValidation{Message: message}
}

// SetError sets the error.
func (e *ErrValidation) SetError(err error) {
	e.Err = err
}
