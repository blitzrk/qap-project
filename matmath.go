package main

import ()

func apply(m [][]float64, f func(float64) float64) [][]float64 {
	res := make([][]float64, len(m))
	for ri, r := range m {
		res[ri] = make([]float64, len(r))
		for ci, c := range r {
			res[ri][ci] = f(c)
		}
	}
	return res
}

func sum(m [][]float64) float64 {
	var s float64
	for _, r := range m {
		for _, c := range r {
			s += c
		}
	}
	return s
}
