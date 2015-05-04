package matrix

import ()

type Element float64

type Matrix [][]Element

func (m Matrix) At(r, c int) Element {
	return m[r][c]
}
