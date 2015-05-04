package main

import (
	"github.com/blitzrk/qap-project/matrix"
	_ "github.com/blitzrk/qap-project/stats"
)

func generate(n int, spread float64) matrix.Matrix {
	return matrix.Matrix(make([][]matrix.Element, 0))
}
