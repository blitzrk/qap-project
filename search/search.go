package search

import (
	"github.com/blitzrk/qap-project/matrix"
	"math"
	"time"
)

type Runner struct {
	NumCPU int
	Cost   matrix.Matrix4D
	fs     *fastStore
}

func (r *Runner) Run(stop <-chan int, resultChan chan<- []uint8) {
	r.fs = NewFS()
	limit := make(chan int, r.NumCPU)
	done := make(chan []uint8)

loop:
	for {
		select {
		case limit <- 1:
			go r.search(done)
		case res := <-done:
			resultChan <- []uint8(res)
			<-limit
		case <-stop:
			break loop
		}
	}
}

func (r *Runner) search(done chan<- []uint8) {
	n := len(r.Cost)
	p := RandPerm(n)

	time.Sleep(1 * time.Second)
	done <- []uint8(p)
}

func greedy(p permutation, r *Runner, done chan<- []uint8) {
	if r.fs.Test(p) {
		done <- []uint8{}
	}
	r.fs.Store(p)

	// Find best neighbor
	var bestPerm permutation
	bestScore := math.Inf(1)
	for _, v := range p.Neighborhood() {
		score := Objective(r.Cost, v)
		if score <= bestScore {
			bestScore = score
			bestPerm = v
		}
	}

	// If best neighbor is worse than current, then we found a local min
	if Objective(r.Cost, p) > bestScore {
		done <- []uint8(p)
		return
	}

	// Otherwise follow the best neighbor
	greedy(bestPerm, r, done)
}
