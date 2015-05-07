package main

import (
	"fmt"
	"github.com/blitzrk/qap-project/data"
	"github.com/blitzrk/qap-project/search"
	"runtime"
	"time"
)

var (
	fact func(int) uint
)

func init() {
	fact = factorial()
}

func main() {
	// Setup data generator
	n := 7
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
	maxTime := time.NewTimer(15 * time.Minute)
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
