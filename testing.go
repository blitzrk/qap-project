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

var (
	fact func(int) uint
)

func AllTests() {
	fact = factorial()
	// testQAPLIBData()
	// testGen()
	// testPermutation()
	testSearch()
	// testHash()
}

func testHash() {
	p1 := search.NewPerm([]uint8{1, 2, 4, 3})
	p2 := search.NewPerm([]uint8{4, 1, 2, 3})
	fmt.Println(p1)
	fmt.Println(p2)
}

func testSearch() {
	// Setup data generator
	n := 6
	gen := data.New(n, 100000)

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
	maxTime := time.NewTimer(5 * time.Minute)
	runner := &search.Runner{
		NumCPU:    runtime.NumCPU(),
		Cost:      cost,
		VarCutoff: 0,
		ProbSize:  fact(n),
	}

	// Run on all 4 cores
	quit := make(chan int)
	results := make(chan *search.Result)
	completed := make(chan bool)
	go runner.Run(quit, results, completed)

loop:
	for {
		select {
		case res := <-results:
			if res != nil {
				fmt.Println(res.Score, res.Perm)
			}
		case <-completed:
			// Bug: may lose last few solutions due to race condition
			fmt.Println("Completed entire search.")
			break loop
		case <-maxTime.C:
			quit <- 1
			fmt.Println("Time out.")
			break loop
		}
	}
}

func testPermutation() {
	fs := search.NewFS(2)
	p1 := search.NewPerm([]uint8{1, 2})
	p2 := search.NewPerm([]uint8{2, 1})
	fs.Store(p1)

	fmt.Println(p1.Hash())
	fmt.Println(p2.Hash())
	fmt.Println(fs.Test(p1), fs.Test(p2))
	fmt.Println(fs.Full())
	fs.Store(p2)
	fmt.Println(fs.Full())
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

func factorial() func(int) uint {
	memo := []uint{1}

	fact := func(i int) uint {
		if i >= len(memo) {
			memo = append(memo, uint(i)*fact(i-1))
		}
		return memo[i]
	}

	return fact
}
