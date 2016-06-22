// Package jchash implements a fast, minimal memory, consistent hash algorithm
// based on the [paper](http://arxiv.org/pdf/1406.2294.pdf) by John Lamping
// and Eric Veach.
package jchash

// JCHash defines a constistent hash range of buckets.
type JCHash struct {
	Buckets int64
}

// Bucket returns the bucket for a given key.
func (h *JCHash) Bucket(key uint64) int64 {
	if h.Buckets <= 0 {
		panic("invalid bucket amount")
	}
	return jumpConsistentHash(key, h.Buckets)
}

func jumpConsistentHash(key uint64, nBuckets int64) int64 {
	var b, j int64 = -1, 0
	for j < nBuckets {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(1<<31) / float64(key>>33+1)))
	}
	return b
}
