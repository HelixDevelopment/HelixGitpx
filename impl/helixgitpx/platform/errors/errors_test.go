package errors_test

import (
	stderrors "errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/helixgitpx/platform/errors"
)

func TestNew_RoundTripFields(t *testing.T) {
	cause := stderrors.New("underlying")
	e := errors.New(codes.NotFound, "repo", "branch %q missing", "main").
		Wrap(cause).
		With("ref", "refs/heads/main")

	if e.Code != codes.NotFound {
		t.Errorf("Code = %v, want NotFound", e.Code)
	}
	if e.Domain != "repo" {
		t.Errorf("Domain = %q, want repo", e.Domain)
	}
	if e.Message != `branch "main" missing` {
		t.Errorf("Message = %q, want quoted branch", e.Message)
	}
	if !stderrors.Is(e, cause) {
		t.Errorf("Is(cause) = false, want true")
	}
	if e.Details["ref"] != "refs/heads/main" {
		t.Errorf("Details[ref] = %v, want refs/heads/main", e.Details["ref"])
	}
}

func TestError_HTTPStatus(t *testing.T) {
	cases := []struct {
		code codes.Code
		want int
	}{
		{codes.OK, http.StatusOK},
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.NotFound, http.StatusNotFound},
		{codes.PermissionDenied, http.StatusForbidden},
		{codes.Unauthenticated, http.StatusUnauthorized},
		{codes.ResourceExhausted, http.StatusTooManyRequests},
		{codes.FailedPrecondition, http.StatusPreconditionFailed},
		{codes.Aborted, http.StatusConflict},
		{codes.Unavailable, http.StatusServiceUnavailable},
		{codes.DeadlineExceeded, http.StatusGatewayTimeout},
		{codes.Unimplemented, http.StatusNotImplemented},
		{codes.Internal, http.StatusInternalServerError},
	}
	for _, c := range cases {
		e := errors.New(c.code, "x", "msg")
		if got := e.HTTPStatus(); got != c.want {
			t.Errorf("code %v: HTTPStatus() = %d, want %d", c.code, got, c.want)
		}
	}
}

func TestError_IsByCode(t *testing.T) {
	e1 := errors.New(codes.NotFound, "repo", "x")
	e2 := errors.New(codes.NotFound, "repo", "y")
	if !stderrors.Is(e1, e2) {
		t.Errorf("Is(same code+domain) = false, want true")
	}
	e3 := errors.New(codes.InvalidArgument, "repo", "z")
	if stderrors.Is(e1, e3) {
		t.Errorf("Is(diff code) = true, want false")
	}
}
