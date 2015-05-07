package main

import (
	"fmt"
	"github.com/blitzrk/qap-project/dat"
	"github.com/blitzrk/qap-project/data"
	"github.com/blitzrk/qap-project/matrix"
	"github.com/blitzrk/qap-project/search"
	"io/ioutil"
	"runtime"
	"time"
)

func AllTests() {
	// 	testQAPLIBData()
	// 	testGen()
	// 	testPermutation()
	testSearch()
	//  testHash()
}

func testHash() {
	p0 := search.NewPerm([]uint8{4, 3, 2, 1})
	p1 := search.NewPerm([]uint8{1, 2, 4, 3})
	p2 := search.NewPerm([]uint8{4, 1, 2, 3})
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println()
	fmt.Println("Ordered by hash from zero:")
	for i := 0; i < 24; i++ {
		fmt.Println(p0.NextNeighbor())
	}
}

func testSearch() {
	// Setup data generator
	gen := data.New(13, 100000)

	// Generate data
	dist, err := gen.Distance()
	if err != nil {
		panic(err)
	}
	flow, err := gen.Flow(1 / 3)
	if err != nil {
		panic(err)
	}
	cost, err := dist.Combine(flow)
	if err != nil {
		panic(err)
	}

	// Setup runner
	maxTime := time.After(5 * time.Minute)
	runner := &search.Runner{
		NumCPU:     runtime.NumCPU(),
		Cost:       cost,
		VarCutoff:  1000 * 1000,
		SampleSize: 1000,
	}

	// Run on all 4 cores
	quit := make(chan int)
	results := make(chan *search.Result)
	go runner.Run(quit, results)

loop:
	for {
		select {
		case res := <-results:
			if res != nil {
				fmt.Println(res.Score, res.Perm)
			}
		case <-maxTime:
			quit <- 1
			fmt.Println("Time out.")
			break loop
		}
	}
}

func testPermutation() {
	fs := search.NewFS()
	p1 := search.NewPerm([]uint8{1, 2, 4, 3})
	p2 := search.NewPerm([]uint8{4, 1, 2, 3})
	fs.Store(p1)

	fmt.Println()
	fmt.Println(p1.Hash())
	fmt.Println()
	fmt.Println(p2.Hash())
	fmt.Println()
	fmt.Println(fs.Test(p1), fs.Test(p2))
}

func testGen() {
	gen := data.New(5, 10000)

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
