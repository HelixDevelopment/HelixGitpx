// Package errors provides the canonical HelixGitpx error type.
//
// Usage:
//
//	err := errors.New(codes.NotFound, "repo", "ref %q missing", name).
//		Wrap(cause).
//		With("ref", ref)
//	return err
//
// HTTP handlers map errors to RFC 7807 problem documents via err.ToProblem.
package errors
