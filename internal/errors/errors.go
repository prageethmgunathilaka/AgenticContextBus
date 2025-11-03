package errors

import (
	"errors"
	"fmt"
)

var (
	ErrExpiredToken = errors.New("token expired")
)

// ErrorCode represents API error codes
type ErrorCode string

const (
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrorCodeValidationError    ErrorCode = "VALIDATION_ERROR"
	ErrorCodeRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// ACBError represents an ACB-specific error
type ACBError struct {
	Code    ErrorCode
	Message string
	Details map[string]string
	Err     error
}

func (e *ACBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ACBError) Unwrap() error {
	return e.Err
}

// NewACBError creates a new ACBError
func NewACBError(code ErrorCode, message string) *ACBError {
	return &ACBError{
		Code:    code,
		Message: message,
		Details: make(map[string]string),
	}
}

// WithDetails adds details to the error
func (e *ACBError) WithDetails(key, value string) *ACBError {
	e.Details[key] = value
	return e
}

// WithError wraps another error
func (e *ACBError) WithError(err error) *ACBError {
	e.Err = err
	return e
}

// Is checks if error matches the code
func Is(err error, code ErrorCode) bool {
	var acbErr *ACBError
	if errors.As(err, &acbErr) {
		return acbErr.Code == code
	}
	return false
}

// Convenience functions for common errors
func Unauthorized(message string) *ACBError {
	return NewACBError(ErrorCodeUnauthorized, message)
}

func Forbidden(message string) *ACBError {
	return NewACBError(ErrorCodeForbidden, message)
}

func NotFound(message string) *ACBError {
	return NewACBError(ErrorCodeNotFound, message)
}

func ValidationError(message string) *ACBError {
	return NewACBError(ErrorCodeValidationError, message)
}

func RateLimitExceeded(message string) *ACBError {
	return NewACBError(ErrorCodeRateLimitExceeded, message)
}

func InternalError(message string) *ACBError {
	return NewACBError(ErrorCodeInternalError, message)
}

func ServiceUnavailable(message string) *ACBError {
	return NewACBError(ErrorCodeServiceUnavailable, message)
}
