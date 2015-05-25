package dat

import (
	"bytes"
	"github.com/blitzrk/qap-project/matrix"
	"strconv"
	"strings"
)

// Reads a .dat file where the entries are matricies of numbers
// separated by empty new lines and individually aligned with
// whitespace between row entries and newlines between rows
func Read(file []byte) []matrix.Matrix {
	data := make([]matrix.Matrix, 0)
	for _, mat := range bytes.Split(file, []byte("\n\n")) {
		temp := make(matrix.Matrix, 0)
		didParse := true
	Element:
		for _, row := range bytes.Split(mat, []byte("\n")) {
			floatsAsStrings := strings.Fields(string(row))
			elms := make([]matrix.Element, 0)
			for _, s := range floatsAsStrings {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					didParse = false
					break Element
				}
				elms = append(elms, matrix.Element(f))
			}
			temp = append(temp, elms)
		}

		if didParse && len(temp[0]) > 0 {
			data = append(data, temp)
		}
	}
	return data
}
