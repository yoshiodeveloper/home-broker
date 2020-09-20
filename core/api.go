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

// SetError sets the error.
func (e *APIError) SetError(err error) {
	e.Err = err
}
