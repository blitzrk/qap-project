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
	fs.Store(p1)

	fmt.Println()
	fmt.Println(p1.Hash())
	fmt.Println()
	fmt.Println(p2.Hash())
	fmt.Println()
	fmt.Println(fs.Test(p1), fs.Test(p2))
}

func testGen() {
	gen := data.New(13, 100000)

	dist, err := gen.Distance()
	if err != nil {
		panic(err)
	}
	flow, err := gen.Flow(1 / 3)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(dist)
	fmt.Println()
	fmt.Println(flow)
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
