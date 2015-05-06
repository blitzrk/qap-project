package matrix

import (
	"errors"
)

var (
	SizeMismatchError = errors.New("Matrix sizes do not match")
	NotSquareError    = errors.New("One or more matrices are not square")
)

type Element float64

type Matrix [][]Element

type Matrix4D [][][][]Element

func (m Matrix) At(r, c int) Element {
	return m[r][c]
}

func (m1 Matrix) Combine(m2 Matrix) (Matrix4D, error) {
	if len(m1) != len(m2) {
		return nil, SizeMismatchError
	}

	n := len(m1)
	m := make(Matrix4D, n)
	for i := 0; i < n; i++ {
		if len(m1[i]) != n || len(m2[i]) != n {
			return nil, NotSquareError
		}
		m[i] = make([][][]Element, n)
		for j := 0; j < n; j++ {
			m[i][j] = make([][]Element, n)
			for x := 0; x < n; x++ {
				m[i][j][x] = make([]Element, n)
				for y := 0; y < n; y++ {
					m[i][j][x][y] = m1[i][j] * m2[x][y]
				}
			}
		}
	}
	return m, nil
}

func (m Matrix4D) At(w, x, y, z int) Element {
	return m[w][x][y][z]
}
