package dat

import (
	"bytes"
	"strconv"
	"strings"
)

// Reads a .dat file where the entries are matricies of numbers
// separated by empty new lines and individually aligned with
// whitespace between row entries and newlines between rows
func Read(file []byte) [][][]float64 {
	data := make([][][]float64, 0)
	for _, matrix := range bytes.Split(file, []byte("\n\n")) {
		temp := make([][]float64, 0)
		didParse := true
	Element:
		for _, row := range bytes.Split(matrix, []byte("\n")) {
			floatsAsStrings := strings.Fields(string(row))
			floats := make([]float64, 0)
			for _, s := range floatsAsStrings {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					didParse = false
					break Element
				}
				floats = append(floats, f)
			}
			temp = append(temp, floats)
		}

		if didParse && len(temp[0]) > 0 {
			data = append(data, temp)
		}
	}
	return data
}
