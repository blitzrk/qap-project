package main

import (
	"math/big"
)

var (
	memo []uint64 = []uint64{1}
)

type Permutation []uint64
type FastStore *big.Int

func NewFS() FastStore {
	return FastStore(&big.Int{})
}

// Hashes a permutation of fixed length n to a number between
// 0 and n!-1 so that a related state may be toggled in a bit
// array.
func (p Permutation) Hash() uint64 {
	return hash(p, 0)
}

func (p Permutation) StoreIn(store FastStore) {
	(*store).SetBit(store, int(p.Hash()), uint(1))
}

func (p Permutation) CheckIn(store FastStore) bool {
	// Bug: Bit takes int, so this only works for permutations
	// up to 20 elements for 64-bit computers
	return (*store).Bit(int(p.Hash())) == uint(1)
}

func hash(p Permutation, pos int) uint64 {
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
