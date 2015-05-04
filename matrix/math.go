package matrix

import ()

func (m Matrix) Apply(f func(Element) Element) Matrix {
	res := make(Matrix, len(m))
	for ri, r := range m {
		res[ri] = make([]Element, len(r))
		for ci, c := range r {
			res[ri][ci] = f(c)
		}
	}
	return res
}

func (m Matrix) Sum() float64 {
	var s float64
	for _, r := range m {
		for _, c := range r {
			s += float64(c)
		}
	}
	return s
}
