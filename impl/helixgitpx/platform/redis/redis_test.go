package redis_test

import (
	"errors"
	"testing"

	hr "github.com/helixgitpx/platform/redis"
)

func TestKey_AppliesNamespace(t *testing.T) {
	c := hr.Client{Namespace: "hello"}
	got := c.Key("greeting", "world")
	want := "hello:greeting:world"
	if got != want {
		t.Errorf("Key = %q, want %q", got, want)
	}
}

func TestIsUnavailable(t *testing.T) {
	if !hr.IsUnavailable(hr.ErrUnavailable) {
		t.Errorf("sentinel not classified")
	}
	if hr.IsUnavailable(errors.New("other")) {
		t.Errorf("other err misclassified")
	}
}
