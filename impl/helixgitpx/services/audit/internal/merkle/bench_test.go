package merkle

import (
	"crypto/rand"
	"testing"
)

func makeLeaves(n, size int) [][]byte {
	leaves := make([][]byte, n)
	for i := range leaves {
		leaves[i] = make([]byte, size)
		_, _ = rand.Read(leaves[i])
	}
	return leaves
}

func BenchmarkRoot_1kLeaves(b *testing.B) {
	leaves := makeLeaves(1000, 256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Root(leaves)
	}
}

func BenchmarkRoot_100kLeaves(b *testing.B) {
	leaves := makeLeaves(100_000, 256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Root(leaves)
	}
}
