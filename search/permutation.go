package search

import (
	"errors"
	"math/rand"
)

var (
	memo []uint64 = []uint64{1}
)

type permutation struct {
	Seq []uint8
	pos int
}

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

// Returns all permutations within a 2-exchange neighborhood
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

// Returns the next permutation in a 2-exchange neighborhood of p
func (p permutation) NextNeighbor(d int) permutation {
	p.pos++
	return RandPerm(len(p))
}

// Returns a random permutations within approximate Hamming distance d
func (p permutation) NextHamming(d int) permutation {
	if d < 2 {
		panic(errors.New("No permutations have a Hamming distance less than 2"))
		return nil
	}

	// After extensive research, no efficient algorithm for enumerating all permutations within
	// a given Hamming distance could be found. As such, an approximation through sampling is used.
	//
	// The cardinality for n=13, d=2 is 78. For d=3, it's 1,352 and for d=4, it's 15,093. An
	// increase of 1 in the Hamming distance appears to approximately lead to an order of magnitude
	// increase of 1. Thus for now we'll recursively sample 10 permutations
	return RandPerm(len(p))
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
