package main

import (
	"fmt"
)

var (
	memo []uint64 = []uint64{1}
)

func main() {
	p1 := Permutation{4, 3, 2, 1}
	p2 := Permutation{3, 1, 2, 4}

	fmt.Println("Hello, playground")

	fmt.Println(p1, p1.Hash())
	fmt.Println(p2, p2.Hash())
}

type Permutation []uint64

func (p Permutation) Hash() uint64 {
	return hash(p, 0)
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

// BigInts shouldn't be needed
//
// func factorial(n *big.Int) (result *big.Int) {
// 	result = new(big.Int)

// 	switch n.Cmp(&big.Int{}) {
// 	case -1, 0:
// 		result.SetInt64(1)
// 	default:
// 		fmt.Println("n = ", n)
// 		result.Set(n)
// 		var one big.Int
// 		one.SetInt64(1)
// 		result.Mul(result, factorial(n.Sub(n, &one)))
// 	}
// 	return
// }
