package acb

import "fmt"

// Error types for SDK
var (
	ErrConnectionFailed = fmt.Errorf("connection failed")
	ErrUnauthorized     = fmt.Errorf("unauthorized")
	ErrNotFound         = fmt.Errorf("not found")
	ErrValidationFailed = fmt.Errorf("validation failed")
)

// SDKError represents an SDK error
type SDKError struct {
	Code    string
	Message string
	Err     error
}

func (e *SDKError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SDKError) Unwrap() error {
	return e.Err
}

// NewSDKError creates a new SDK error
func NewSDKError(code, message string) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
	}
}

// WithError wraps another error
func (e *SDKError) WithError(err error) *SDKError {
	e.Err = err
	return e
}

