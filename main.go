package main

import (
	"fmt"
	"github.com/blitzrk/qap-project/dat"
	"github.com/blitzrk/qap-project/matrix"
	"io/ioutil"
)

func main() {
	readQAPLIBData()
}

func readQAPLIBData() {
	data, err := readDat("bur26a.dat")
	if err != nil {
		fmt.Println(err)
		return
	}

	n := data[0][0][0]
	times := data[1]
	freqs := data[2]
	totalF := freqs.Sum()

	fmt.Println(n)
	fmt.Println()
	fmt.Println(times)
	fmt.Println()
	fmt.Println(freqs)
	fmt.Println()
	fmt.Println(totalF)
}

func readDat(fname string) ([]matrix.Matrix, error) {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	return dat.Read(f), nil
}
