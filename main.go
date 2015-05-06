package main

import (
	"fmt"
	"github.com/blitzrk/qap-project/dat"
	"github.com/blitzrk/qap-project/data"
	"github.com/blitzrk/qap-project/matrix"
	"github.com/blitzrk/qap-project/search"
	"io/ioutil"
)

func main() {
	testQAPLIBData()
	testGen()
	testPermutation()
}

func testPermutation() {
	fs := search.NewFS()
	p1 := search.Permutation{1, 2, 4, 3}
	p2 := search.Permutation{4, 1, 2, 3}
	p1.StoreIn(fs)

	fmt.Println()
	fmt.Println(p1.Hash())
	fmt.Println()
	fmt.Println(p2.Hash())
	fmt.Println()
	fmt.Println(p1.CheckIn(fs), p2.CheckIn(fs))
}

func testGen() {
	gen := data.New(13, 100000)
	graph, err := gen.Flow(1 / 3)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(graph)
}

func testQAPLIBData() {
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
