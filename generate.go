package main

import (
	"fmt"
	"github.com/blitzrk/qap-project/matrix"
	"math"
	"math/rand"
)

func genData(n int, spread float64) (matrix.Matrix, error) {
	if spread < 0 || spread >= 1 {
		return nil, fmt.Errorf("Error: spread must be between 0 and 1.")
	}

	// Create Zipf generator
	scale := 1000
	r := rand.New(rand.NewSource(0))
	z := rand.NewZipf(r, 1.01, float64(n), uint64(scale))

	// Populate frequencies of unigrams
	k := make([]float64, n)
	for i := 0; i < n; i++ {
		k[i] = float64(z.Uint64())
	}
	fmt.Println()
	fmt.Println(k)

	// Populate ideal bigram matrix
	g := matrix.Matrix(make([][]matrix.Element, n))
	for i := 0; i < n; i++ {
		g[i] = make([]matrix.Element, n)
		for j := 0; j < n; j++ {
			e := (rand.Float64() - 0.5) * spread
			g[i][j] = matrix.Element((k[i] * k[j]) * (1 + e))
		}
	}

	// Scale back to 100,000 total freq
	totalF := g.Sum()
	s := 100000 / totalF
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			g[i][j] = matrix.Element(math.Floor(float64(g[i][j]) * s))
		}
	}

	return g, nil
}
