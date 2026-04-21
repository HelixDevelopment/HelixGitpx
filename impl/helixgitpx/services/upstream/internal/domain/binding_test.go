package domain

import (
	"errors"
	"testing"
)

func TestDirection_Allows(t *testing.T) {
	cases := []struct {
		dir   Direction
		op    Op
		allow bool
	}{
		{DirectionReadOnly, OpRead, true},
		{DirectionReadOnly, OpWrite, false},
		{DirectionReadOnly, OpReceiveWebhook, false},
		{DirectionWrite, OpRead, true},
		{DirectionWrite, OpWrite, true},
		{DirectionWrite, OpReceiveWebhook, true},
		{DirectionBidirectional, OpRead, true},
		{DirectionBidirectional, OpWrite, true},
		{DirectionBidirectional, OpReceiveWebhook, true},
	}
	for _, c := range cases {
		if got := c.dir.Allows(c.op); got != c.allow {
			t.Errorf("dir=%d op=%d want %v got %v", c.dir, c.op, c.allow, got)
		}
	}
}

func TestValidate(t *testing.T) {
	base := BindingInput{RepoID: "r1", Provider: "github", RawURL: "https://github.com/o/r.git", Direction: DirectionWrite}
	if err := Validate(base); err != nil {
		t.Fatalf("happy path rejected: %v", err)
	}

	tests := map[string]struct {
		mut  func(*BindingInput)
		want error
	}{
		"empty repo":      {func(b *BindingInput) { b.RepoID = "" }, ErrEmptyRepoID},
		"unknown provider": {func(b *BindingInput) { b.Provider = "phabricator" }, ErrUnknownProvider},
		"bad URL":         {func(b *BindingInput) { b.RawURL = "not a url" }, ErrInvalidURL},
		"ftp scheme":      {func(b *BindingInput) { b.RawURL = "ftp://foo/bar.git" }, ErrInvalidURL},
		"no direction":    {func(b *BindingInput) { b.Direction = DirectionUnspecified }, ErrInvalidDirection},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			in := base
			tc.mut(&in)
			if got := Validate(in); !errors.Is(got, tc.want) {
				t.Fatalf("want %v got %v", tc.want, got)
			}
		})
	}
}

func TestProviders_AllTwelve(t *testing.T) {
	if len(Providers) != 12 {
		t.Fatalf("expected 12 providers at GA, got %d", len(Providers))
	}
}
