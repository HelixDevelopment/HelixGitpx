package merkle_test

import (
	"testing"

	"github.com/helixgitpx/helixgitpx/services/audit/internal/merkle"
)

func TestRoot_SingleLeaf(t *testing.T) {
	root := merkle.Root([][]byte{[]byte("a")})
	if len(root) == 0 {
		t.Fatal("empty root")
	}
}

func TestRoot_TwoLeaves_Deterministic(t *testing.T) {
	r1 := merkle.Root([][]byte{[]byte("a"), []byte("b")})
	r2 := merkle.Root([][]byte{[]byte("a"), []byte("b")})
	if string(r1) != string(r2) {
		t.Error("root not deterministic")
	}
	r3 := merkle.Root([][]byte{[]byte("b"), []byte("a")})
	if string(r1) == string(r3) {
		t.Error("root must be order-sensitive")
	}
}

func TestRoot_Empty(t *testing.T) {
	if r := merkle.Root(nil); r != nil {
		t.Errorf("empty input = %v, want nil", r)
	}
}
