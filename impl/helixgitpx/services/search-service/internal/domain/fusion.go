// Package domain encodes search-service algorithms: the Reciprocal Rank
// Fusion score combiner plus hit-sanitisation invariants.
package domain

import (
	"sort"
)

// RRFk is the conventional constant used in Reciprocal Rank Fusion; 60 is
// the value from the original Cormack et al. paper.
const RRFk = 60

// Ranking is a single engine's ordered hit list. The position in the slice
// is the rank (0 = top).
type Ranking struct {
	Engine string    // "meilisearch" | "qdrant" | "zoekt"
	Hits   []string  // hit IDs in descending relevance order
	Weight float64   // multiplier applied to 1/(k+rank); defaults to 1.0
}

// FusedHit is the merged representation returned to the client.
type FusedHit struct {
	ID    string
	Score float64
	// PerEngine maps engine name → rank contribution (1/(k+rank) * weight).
	PerEngine map[string]float64
}

// Fuse applies Reciprocal Rank Fusion across arbitrary many rankings and
// returns hits sorted descending by fused score. An input with zero Weight
// is treated as Weight=1.0.
func Fuse(rankings []Ranking) []FusedHit {
	scores := map[string]*FusedHit{}
	for _, r := range rankings {
		w := r.Weight
		if w == 0 {
			w = 1.0
		}
		for rank, id := range r.Hits {
			contrib := w / float64(RRFk+rank+1)
			if hit, ok := scores[id]; ok {
				hit.Score += contrib
				hit.PerEngine[r.Engine] = contrib
			} else {
				scores[id] = &FusedHit{
					ID:        id,
					Score:     contrib,
					PerEngine: map[string]float64{r.Engine: contrib},
				}
			}
		}
	}
	out := make([]FusedHit, 0, len(scores))
	for _, h := range scores {
		out = append(out, *h)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// TopK trims a fused list to the first K results. A K ≤ 0 returns everything.
func TopK(hits []FusedHit, k int) []FusedHit {
	if k <= 0 || k >= len(hits) {
		return hits
	}
	return hits[:k]
}
