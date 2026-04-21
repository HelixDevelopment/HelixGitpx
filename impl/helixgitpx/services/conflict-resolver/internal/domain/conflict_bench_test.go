package domain

import "testing"

func BenchmarkClassify(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Classify(true, false, false, false)
		_ = Classify(false, false, true, false)
		_ = Classify(false, true, false, false)
		_ = Classify(false, false, false, true)
	}
}

func BenchmarkTransition(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Transition(StatusOpen, StatusProposed)
		_ = Transition(StatusProposed, StatusResolved)
	}
}
