package search

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

var (
	memo []uint64 = []uint64{1}
)

type permutation struct {
	Seq []uint8
	pos int
}

func (p *permutation) String() string {
   return fmt.Sprint(p.pos, ": ", p.Seq)
}

// Create a permutation of 1...n from an int slice
func NewPerm(p []uint8) *permutation {
	h := hash(p, 0)
	return &permutation{p, int(h)}
}

// Create random permutation of 1...n
func RandPerm(n int) *permutation {
	p := rand.Perm(n)
	r := make([]uint8, len(p))

	for i, v := range p {
		r[i] = uint8(v + 1)
	}

	return NewPerm(r)
}

// DEPRECATED:
// Returns all permutations within a 2-exchange neighborhood
func (p *permutation) Neighborhood() []*permutation {
	n := len(p.Seq)
	perms := make([]*permutation, 0, n*(n-1)/2)

	// Find 2-exchange neighborhood
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			perm := make([]uint8, len(p.Seq))
			copy(perm, p.Seq)
			perm[j], perm[i] = p.Seq[i], p.Seq[j]
			perms = append(perms, NewPerm(perm))
		}
	}

	return perms
}

// Returns the next permutation in a 2-exchange neighborhood of p
func (p *permutation) NextNeighbor() *permutation {
  // Cycle position 1
  p.pos++
  if p.pos == int(fact(uint64(len(p.Seq)))) {
    p.pos = 0
  }
  
  // Unhash by position
	s := p.Unhash()

	return NewPerm(s)
}

// After extensive research, no efficient algorithm for enumerating all permutations within
// a given Hamming distance could be found. As such, an approximation through sampling is used.
//
// IRRELEVANT: The cardinality for n=13, d=2 is 78. For d=3, it's 1,352 and for d=4, it's 15,093.
// An increase of 1 in the Hamming distance appears to approximately lead to an order of magnitude
// increase of 1 near n=13.
//
// Returns a random permutations within approximate Hamming distance d
func (p *permutation) NextHamming(d int) *permutation {
	if d < 2 {
		panic(errors.New("No permutations have a Hamming distance less than 2"))
		return nil
	}

	s := make([]uint8, len(p.Seq))
	copy(s, p.Seq)

	for d >= 2 {
		if d == 3 {
			s = Exchange3Rand(s)
			d -= 3
		} else {
			s = Exchange2Rand(s)
			d -= 2
		}
	}

	return NewPerm(s)
}

func Exchange2Rand(seq []uint8) []uint8 {
	n := len(seq)
	s := make([]uint8, n)
	copy(s, seq)
	random := rand.Perm(n)

	s[random[0]], s[random[1]] = s[random[1]], s[random[0]]
	return s
}

func Exchange3Rand(seq []uint8) []uint8 {
	n := len(seq)
	s := make([]uint8, n)
	copy(s, seq)
	random := rand.Perm(n)

	s[random[0]], s[random[1]] = s[random[1]], s[random[0]]
	s[random[1]], s[random[2]] = s[random[2]], s[random[1]]
	return s
}

// Hashes a permutation of fixed length n to a number between
// 0 and n!-1 so that a related state may be toggled in a bit
// array.
func (p *permutation) Hash() uint64 {
	return hash(p.Seq, 0)
}

func (p *permutation) Unhash() []uint8 {
  n := len(p.Seq)
	s := make([]uint8, n)
	
	ints := make([]uint8, n)
	for i, _ := range ints {
	  ints[i] = uint8(i+1)
	}
	
	hsh := float64(p.pos)
	for i := 0; i < n; i++ {
	  fac := fact(uint64(n - 1 - i))
	  order := int(math.Floor(hsh / float64(fac)))
	  s[i] = ints[order]
	  
	  // Keep track of the remaining elements (in order)
	  ints = append(ints[:order], ints[order+1:]...)
	  
	  // Subtract out contribution to hash
	  factor := s[i]
  	for j := 0; j < i; j++ {
  		if s[j] < s[i] {
  			factor--
  		}
  	}
  	factor--
  	
	  hsh -= float64(factor) * float64(fac)
	}
	
	return s
}

func hash(seq []uint8, pos int) uint64 {
	n := len(seq)
	if pos >= n {
		return 0
	}

	s, order := seq[pos], seq[pos]
	for i := 0; i < pos; i++ {
		if seq[i] < s {
			order--
		}
	}

	return uint64(order-1)*fact(uint64(n-1-pos)) + hash(seq, pos+1)
}

func fact(i uint64) uint64 {
	if i >= uint64(len(memo)) {
		memo = append(memo, i*fact(i-1))
	}
	return memo[i]
}
