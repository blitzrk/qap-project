package search

import (
	"github.com/blitzrk/qap-project/matrix"
	"math"
	"math/rand"
	"runtime"
)

type Runner struct {
	NumCPU      int
	Cost        matrix.Matrix4D
	StartRadius int
	SampleSize  int
	fs          *fastStore
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
			go r.search(p, done)
		case res := <-done:
			resultChan <- []uint8(res)
			<-limit
		case <-stop:
			break loop
		}
	}
}

type runResult struct {
	Perm   permutation
	Score  float64
	Nils   int
	Var    float64
	Center permutation
	FinalR int
}

func (r *Runner) search(perm permutation, done chan<- []uint8) {
	collect := make(chan *runResult)
	go r.findBestNeighbor(perm, collect)

	results := make([]*runResult, 0, r.Sample)
	for i := 0; i < r.Sample; i++ {
		res := <-collect
		results = append(results, res)
	}

	// Change what gets sent here
	go r.interpret(results, done)
}

func (r *Runner) greedy(p permutation, done chan<- *runResult) {
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
func (r *Runner) findBestNeighbor(center permutation, done chan<- *runResult) {
	n := len(center)
	size := n * (n - 1) / 2
	if size > r.SampleSize {
		size = r.SampleSize
	}

	var bestPerm permutation
	bestScore := math.Inf(1)
	var nils int
	scores := make([]float64, size)

	for i := 0; i < size; i++ {
		neighbor := p.NextNeighbor()
		if v == nil {
			nils++
			continue
		}

		score := r.Objective(neighbor)
		scores[i] = score

		if score <= bestScore {
			bestScore = score
			bestPerm = neighbor
		}
	}

	vari := variance(scores)

	done <- &runResult{
		Perm:   bestPerm,
		Score:  bestScore,
		Nils:   nils,
		Var:    vari,
		Center: center,
		FinalR: 2,
	}
}

// Find best permutation from sampled APPROXIMATE Hamming space
// TODO: predict size of Hamming for max sample size
func (r *Runner) findBestHamming(center permutation, dist int, done chan<- *runResult) {
	var bestPerm permutation
	bestScore := math.Inf(1)
	var nils int
	scores := make([]float64, r.SampleSize)

	for i := 0; i < r.SampleSize; i++ {
		neighbor := p.NextHamming(dist)
		if v == nil {
			nils++
			continue
		}

		score := r.Objective(neighbor)
		scores[i] = score

		if score <= bestScore {
			bestScore = score
			bestPerm = neighbor
		}
	}

	vari := variance(scores)

	done <- &runResult{
		Perm:   bestPerm,
		Score:  bestScore,
		Nils:   nils,
		Var:    vari,
		Center: center,
		FinalR: dist,
	}
}

// Consider the number of nils (followed old path) and
func (r *Runner) interpret(rs *runResult) (permutation, bool) {

	// TODO: Find variance of scores

	bestPerm, _ := findBest(rs, cost)

	return bestPerm, true
}

func variance(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum += v
	}
	mean := sum / len(x)

	var sumsq float64
	for _, v := range x {
		sumsq += math.Pow(v-mean, 2)
	}
	vari := sumsq / len(x)

	return vari
}
