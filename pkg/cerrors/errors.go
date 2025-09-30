// Package cerrors provides custom error handling utilities.
package cerrors

import (
	"github.com/cockroachdb/errors"
)

// New creates a new error with the given message.
func New(msg string) error {
	return errors.New(msg)
}

// Errorf creates a new error with the given format and arguments.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Wrap wraps an error with additional context.
func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// Wrapf wraps an error with additional formatted context.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// WithMessage adds a message to an error without wrapping it.
func WithMessage(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

// WithMessagef adds a formatted message to an error without wrapping it.
func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

// WithDetail adds detail information to an error.
func WithDetail(err error, msg string) error {
	return errors.WithDetail(err, msg)
}

// WithHint adds a hint to an error.
func WithHint(err error, msg string) error {
	return errors.WithHint(err, msg)
}

// Is checks if an error matches a target error.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As checks if an error can be assigned to a target.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the underlying error.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Cause returns the underlying cause of an error.
func Cause(err error) error {
	return errors.Cause(err)
}

// Common application errors
var (
	// ErrNotFound indicates that a resource was not found.
	ErrNotFound = New("resource not found")

	// ErrInvalidInput indicates that input validation failed.
	ErrInvalidInput = New("invalid input")

	// ErrUnauthorized indicates that authentication failed.
	ErrUnauthorized = New("unauthorized")

	// ErrForbidden indicates that authorization failed.
	ErrForbidden = New("forbidden")

	// ErrConflict indicates that a conflict occurred.
	ErrConflict = New("conflict")

	// ErrInternalServer indicates an internal server error.
	ErrInternalServer = New("internal server error")

	// ErrServiceUnavailable indicates that the service is unavailable.
	ErrServiceUnavailable = New("service unavailable")

	// ErrTimeout indicates that an operation timed out.
	ErrTimeout = New("operation timed out")

	// ErrNetworkError indicates a network-related error.
	ErrNetworkError = New("network error")
)

// ErrorCode represents different types of errors
type ErrorCode string

// Error codes
const (
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeInvalidInput       ErrorCode = "INVALID_INPUT"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeForbidden          ErrorCode = "FORBIDDEN"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeInternalServer     ErrorCode = "INTERNAL_SERVER"
	CodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	CodeTimeout            ErrorCode = "TIMEOUT"
	CodeNetworkError       ErrorCode = "NETWORK_ERROR"
)

// AppError represents an application-specific error with a code.
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError.
func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsAppError checks if an error is an AppError with the given code.
func IsAppError(err error, code ErrorCode) bool {
	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// GetErrorCode extracts the error code from an AppError.
func GetErrorCode(err error) ErrorCode {
	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code
	}
	return ""
}