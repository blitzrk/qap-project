package search

import (
	"math/big"
)

type fastStore struct {
	big.Int
	size uint
}

func NewFS(size uint) *fastStore {
	return &fastStore{big.Int{}, size}
}

func (fs *fastStore) Store(p *permutation) {
	fs.Int.SetBit(&fs.Int, int(p.hash), 1)
}

func (fs *fastStore) Test(p *permutation) bool {
	// Bug: Bit takes int, so this only works for permutations
	// up to 20 elements for 64-bit computers
	return fs.Int.Bit(int(p.hash)) == uint(1)
}

func (fs *fastStore) Full() bool {
	full := big.NewInt(0)
	for i := 0; i < int(fs.size); i++ {
		full.SetBit(full, i, 1)
	}
	return fs.Int.Cmp(full) == 0
}
