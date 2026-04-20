// Package errors defines the canonical HelixGitpx error type.
//
// An Error carries a gRPC status code, a domain tag, a human-readable message,
// an optional cause, and a map of structured details. It implements the Go
// error interface, the standard errors.Is/As contract, and supplies an HTTP
// status mapping per RFC 7807.
package errors

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// Error is HelixGitpx's typed error.
type Error struct {
	Code    codes.Code
	Domain  string
	Message string
	Cause   error
	Details map[string]any
}

// New constructs an Error. The message is formatted with fmt.Sprintf if args are supplied.
func New(code codes.Code, domain, format string, args ...any) *Error {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return &Error{Code: code, Domain: domain, Message: msg}
}

// Wrap attaches a cause and returns the receiver.
func (e *Error) Wrap(cause error) *Error {
	e.Cause = cause
	return e
}

// With adds a structured detail and returns the receiver.
func (e *Error) With(key string, value any) *Error {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s/%s] %s: %v", e.Domain, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s/%s] %s", e.Domain, e.Code, e.Message)
}

// Unwrap supports errors.Is/As.
func (e *Error) Unwrap() error { return e.Cause }

// Is returns true when target is an *Error with the same Code and Domain, or
// when the cause matches target.
func (e *Error) Is(target error) bool {
	var other *Error
	if errors.As(target, &other) {
		return e.Code == other.Code && e.Domain == other.Domain
	}
	return false
}

// HTTPStatus maps the gRPC code to the closest HTTP status per RFC 7807.
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.InvalidArgument, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.DataLoss, codes.Internal, codes.Unknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
