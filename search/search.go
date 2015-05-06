package search

import (
	"github.com/blitzrk/qap-project/matrix"
	"math"
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
	Perm   *permutation
	Score  float64
	Nils   int
	Var    float64
	Center *permutation
	FinalR int
}

func (r *Runner) search(perm *permutation, done chan<- []uint8) {
	collect := make(chan *runResult)
	go r.findBestNeighbor(perm, collect)

	// Change what gets sent here
	result := <-collect
	go r.interpret(result, done)
}

func (r *Runner) greedy(p *permutation, done chan<- *runResult) {
	if r.fs.Test(p) {
		done <- nil
		return
	}
	r.fs.Store(p)
}

// Find best permutation
func (r *Runner) findBestNeighbor(center *permutation, done chan<- *runResult) {
	n := len(center.Seq)
	size := n * (n - 1) / 2
	if size > r.SampleSize {
		size = r.SampleSize
	}

	var bestPerm *permutation
	bestScore := math.Inf(1)
	var nils int
	scores := make([]float64, size)

	for i := 0; i < size; i++ {
		neighbor := center.NextNeighbor()
		if neighbor == nil {
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
func (r *Runner) findBestHamming(center *permutation, dist int, done chan<- *runResult) {
	var bestPerm *permutation
	bestScore := math.Inf(1)
	var nils int
	scores := make([]float64, r.SampleSize)

	for i := 0; i < r.SampleSize; i++ {
		neighbor := center.NextHamming(dist)
		if neighbor == nil {
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

// Use a greedy algorithm search for local mins, but also use stats (variance,
// num time ended up on same path) to determine if to expand the search to a
// greater radius (Hamming distance)
func (r *Runner) interpret(rs *runResult, done chan<- []uint8) {
	// TODO: ALGORITHM HERE
	done <- rs.Perm.Seq
}

func variance(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum += v
	}
	mean := sum / float64(len(x))

	var sumsq float64
	for _, v := range x {
		sumsq += math.Pow(v-mean, 2)
	}
	vari := sumsq / float64(len(x))

	return vari
}
