package search

import (
	"math/big"
	"math/rand"
)

var (
	memo []uint64 = []uint64{1}
)

type permutation []uint8

type fastStore struct {
	big.Int
}

func NewFS() *fastStore {
	return &fastStore{big.Int{}}
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

// Hashes a permutation of fixed length n to a number between
// 0 and n!-1 so that a related state may be toggled in a bit
// array.
func (p permutation) Hash() uint64 {
	return hash(p, 0)
}

func (fs *fastStore) Store(p permutation) {
	(*fs).SetBit(&(fs.Int), int(p.Hash()), uint(1))
}

func (fs *fastStore) Test(p permutation) bool {
	// Bug: Bit takes int, so this only works for permutations
	// up to 20 elements for 64-bit computers
	return (*fs).Bit(int(p.Hash())) == uint(1)
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
