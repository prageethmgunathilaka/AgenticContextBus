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
	ErrorCodeInvalidRequest    ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrorCodeConflict          ErrorCode = "CONFLICT"
	ErrorCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeInternalError     ErrorCode = "INTERNAL_ERROR"
)

// APIError represents an API error
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to an API error
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) error {
	if message == "" {
		message = "authentication required"
	}
	return NewAPIError(ErrorCodeUnauthorized, message)
}

// Forbidden creates a forbidden error
func Forbidden(message string) error {
	if message == "" {
		message = "insufficient permissions"
	}
	return NewAPIError(ErrorCodeForbidden, message)
}

// NotFound creates a not found error
func NotFound(message string) error {
	if message == "" {
		message = "resource not found"
	}
	return NewAPIError(ErrorCodeNotFound, message)
}

// Conflict creates a conflict error
func Conflict(message string) error {
	if message == "" {
		message = "resource conflict"
	}
	return NewAPIError(ErrorCodeConflict, message)
}

