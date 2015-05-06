package search

import (
	"math/big"
)

type fastStore struct {
	big.Int
}

func NewFS() *fastStore {
	return &fastStore{big.Int{}}
}

func (fs *fastStore) Store(p *permutation) {
	(*fs).SetBit(&(fs.Int), int(p.Hash()), uint(1))
}

func (fs *fastStore) Test(p *permutation) bool {
	// Bug: Bit takes int, so this only works for permutations
	// up to 20 elements for 64-bit computers
	return (*fs).Bit(int(p.Hash())) == uint(1)
}
