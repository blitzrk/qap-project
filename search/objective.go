package search

import (
	"github.com/blitzrk/qap-project/matrix"
)

func Objective(c matrix.Matrix4D, p permutation) float64 {
	var sum float64
	n := len(p)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum += float64(c.At(i, j, int(p[i])-1, int(p[j])-1))
		}
	}

	return sum
}
