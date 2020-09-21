package core

// APIError represents holds an API error.
type APIError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAPIError creates a new APIError.
func NewAPIError(message string, statusCode int) APIError {
	return APIError{Message: message, StatusCode: statusCode}
}

// NewAPIErrorFromErrValidation a new APIErro from an ErrValidation.
func NewAPIErrorFromErrValidation(errValidation ErrValidation) APIError {
	return APIError{Message: errValidation.Message, StatusCode: 400}
}

// SetError sets the error.
func (e *APIError) SetError(err error) {
	e.Err = err
}
