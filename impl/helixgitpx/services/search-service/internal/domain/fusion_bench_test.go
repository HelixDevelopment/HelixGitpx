package domain

import (
	"fmt"
	"testing"
)

func BenchmarkFuse_ThreeEngines_100Hits(b *testing.B) {
	hits := func(prefix string, n int) []string {
		out := make([]string, n)
		for i := 0; i < n; i++ {
			out[i] = fmt.Sprintf("%s-%d", prefix, i)
		}
		return out
	}
	r1 := Ranking{Engine: "meili", Hits: hits("m", 100), Weight: 1.0}
	r2 := Ranking{Engine: "qdrant", Hits: hits("q", 100), Weight: 1.0}
	r3 := Ranking{Engine: "zoekt", Hits: hits("z", 100), Weight: 2.0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Fuse([]Ranking{r1, r2, r3})
	}
}

func BenchmarkTopK_1000Fused(b *testing.B) {
	fused := make([]FusedHit, 1000)
	for i := range fused {
		fused[i] = FusedHit{ID: fmt.Sprintf("id-%d", i), Score: float64(i)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TopK(fused, 20)
	}
}
