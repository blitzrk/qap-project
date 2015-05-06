package search

import (
	"math/rand"
)

var (
	memo []uint64 = []uint64{1}
)

type permutation []uint8

// Create a permutation of 1...n from an int slice
func NewPerm(p []int) permutation {
	perm := make([]uint8, len(p))
	for i, v := range p {
		perm[i] = uint8(v)
	}
	return perm
}

// Create random permutation of 1...n
func RandPerm(n int) permutation {
	p := rand.Perm(n)

	for i, v := range p {
		p[i] = v + 1
	}

	return NewPerm(p)
}

func (p permutation) Neighborhood() []permutation {
	n := len(p)
	perms := make([]permutation, 0, n*(n-1)/2)

	// Find 2-exchange neighborhood
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			perm := make(permutation, len(p))
			copy(perm, p)
			perm[j], perm[i] = p[i], p[j]
			perms = append(perms, perm)
		}
	}

	return perms
}

// Hashes a permutation of fixed length n to a number between
// 0 and n!-1 so that a related state may be toggled in a bit
// array.
func (p permutation) Hash() uint64 {
	return hash(p, 0)
}

func hash(p permutation, pos int) uint64 {
	n := len(p)
	if pos >= n {
		return 0
	}

	s, factor := p[pos], p[pos]
	for i := 0; i < pos; i++ {
		if p[i] < s {
			factor--
		}
	}

	return uint64(factor-1)*fact(uint64(n-1-pos)) + hash(p, pos+1)
}

func fact(i uint64) uint64 {
	if i >= uint64(len(memo)) {
		memo = append(memo, i*fact(i-1))
	}
	return memo[i]
}
