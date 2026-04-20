// Package merkle builds a SHA-256 binary Merkle tree over ordered leaves.
// Used by audit-service to anchor an hour's events into audit.anchors.
package merkle

import "crypto/sha256"

// Root hashes leaves pairwise until a single root remains. Odd leaves at
// any level are promoted unchanged. Deterministic and order-sensitive.
func Root(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return nil
	}
	level := make([][]byte, len(leaves))
	for i, l := range leaves {
		h := sha256.Sum256(l)
		level[i] = h[:]
	}
	for len(level) > 1 {
		var next [][]byte
		for i := 0; i < len(level); i += 2 {
			if i+1 == len(level) {
				next = append(next, level[i])
				continue
			}
			joined := append(append([]byte{}, level[i]...), level[i+1]...)
			h := sha256.Sum256(joined)
			next = append(next, h[:])
		}
		level = next
	}
	return level[0]
}
