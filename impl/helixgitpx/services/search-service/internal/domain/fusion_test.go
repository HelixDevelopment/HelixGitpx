package domain

import "testing"

func TestFuse_SimplyBlends(t *testing.T) {
	r1 := Ranking{Engine: "meili", Hits: []string{"a", "b", "c"}}
	r2 := Ranking{Engine: "qdrant", Hits: []string{"b", "c", "d"}}

	fused := Fuse([]Ranking{r1, r2})
	if len(fused) != 4 {
		t.Fatalf("want 4 unique hits, got %d", len(fused))
	}
	if fused[0].ID != "b" {
		t.Fatalf("b appears in both lists at top and should rank first; got %s", fused[0].ID)
	}
}

func TestFuse_RespectsWeights(t *testing.T) {
	heavy := Ranking{Engine: "zoekt", Hits: []string{"z1", "a"}, Weight: 5.0}
	light := Ranking{Engine: "meili", Hits: []string{"a", "z1"}, Weight: 1.0}

	fused := Fuse([]Ranking{heavy, light})
	// zoekt ranks z1 first at weight 5, meili ranks a first at weight 1.
	// z1 contribution ≈ 5/(60+1) ≈ 0.082; a contribution ≈ 1/(60+1) + 5/(60+2) ≈ 0.0168+0.0806
	// so z1 should still come out on top.
	if fused[0].ID != "z1" {
		t.Fatalf("weighted fusion should put z1 first, got %s", fused[0].ID)
	}
}

func TestFuse_EmptyInput(t *testing.T) {
	if got := Fuse(nil); len(got) != 0 {
		t.Fatal("empty input must produce empty output")
	}
}

func TestTopK(t *testing.T) {
	hits := []FusedHit{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	if got := TopK(hits, 2); len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
	if got := TopK(hits, 0); len(got) != 3 {
		t.Fatal("k=0 should return all")
	}
	if got := TopK(hits, 99); len(got) != 3 {
		t.Fatal("k >= len should return all")
	}
}

func TestFuse_StableTieBreak(t *testing.T) {
	r1 := Ranking{Engine: "e1", Hits: []string{"x", "y"}}
	// Reverse order gives identical fused scores; result must be alphabetical.
	r2 := Ranking{Engine: "e2", Hits: []string{"y", "x"}}
	fused := Fuse([]Ranking{r1, r2})
	if fused[0].ID != "x" || fused[1].ID != "y" {
		t.Fatalf("want alphabetical on tie, got %v", fused)
	}
}
