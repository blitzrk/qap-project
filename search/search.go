package search

import (
	"github.com/blitzrk/qap-project/matrix"
	"math"
	"math/rand"
	"runtime"
)

type Runner struct {
	NumCPU int
	Cost   matrix.Matrix4D
	Radius int
	Sample int
	fs     *fastStore
}

func (r *Runner) Run(stop <-chan int, resultChan chan<- []uint8) {
	// maximize CPU usage
	runtime.GOMAXPROCS(r.NumCPU)

	r.fs = NewFS()
	n := len(r.Cost)
	limit := make(chan int, r.NumCPU)
	done := make(chan []uint8)

loop:
	for {
		select {
		case limit <- 1:
			p := RandPerm(n)
			go r.search(p, r.Radius, done)
		case res := <-done:
			resultChan <- []uint8(res)
			<-limit
		case <-stop:
			break loop
		}
	}
}

func (r *Runner) search(perm permutation, radius int, done chan<- []uint8) {
	// Radius cannot be less than 2!!!
	h := perm.Hamming(radius)
	collect := make(chan permutation)

	for i := 0; i < r.Sample; i++ {
		p := h[rand.Intn(len(h))]
		go greedy(p, r, collect)
	}

	results := make([]permutation, 0, r.Sample)
	for i := 0; i < r.Sample; i++ {
		res := <-collect
		results = append(results, res)
	}

	// Change what gets sent here
	if v, ok := interpret(results, perm, radius, r.Cost); ok {
		done <- []uint8(v)
	}
}

func greedy(p permutation, r *Runner, done chan<- permutation) {
	if r.fs.Test(p) {
		done <- nil
		return
	}
	r.fs.Store(p)

	bestPerm, bestScore := findBest(p.Neighborhood(), r.Cost)

	// If best neighbor is worse than current, then we found a local min
	if Objective(r.Cost, p) > bestScore {
		done <- p
		return
	}

	// Otherwise follow the best neighbor
	greedy(bestPerm, r, done)
	return
}

// Find best permutation
func findBest(ps []permutation, cost matrix.Matrix4D) (bestPerm permutation, bestScore float64) {

	bestScore = math.Inf(1)
	for _, v := range ps {
		if v == nil {
			continue
		}
		score := Objective(cost, v)
		if score <= bestScore {
			bestScore = score
			bestPerm = v
		}
	}

	return
}

// Consider the number of nils (followed old path) and
func interpret(rs []permutation, perm permutation, radius int, cost matrix.Matrix4D) (permutation, bool) {
	// Count nils
	var nils int
	for _, v := range rs {
		if v == nil {
			nils++
			continue
		}
	}

	// TODO: Find variance of scores

	bestPerm, _ := findBest(rs, cost)

	return bestPerm, true
}
