package search

import (
	"github.com/blitzrk/qap-project/matrix"
	"log"
	"math"
	"os"
	"runtime"
)

var (
	logger *log.Logger
)

func init() {
	file, err := os.OpenFile("go.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}
	logger = log.New(file, "logger: ", log.Lshortfile)
}

type Runner struct {
	NumCPU    int
	Cost      matrix.Matrix4D
	VarCutoff float64
	ProbSize  uint
	fs        *fastStore
}

func (r *Runner) Run(stop <-chan int, resultChan chan<- *Result, complete chan<- bool) {
	// maximize CPU usage
	runtime.GOMAXPROCS(r.NumCPU)
	limit := make(chan int, r.NumCPU)

	r.fs = NewFS(r.ProbSize)
	n := len(r.Cost)
	done := make(chan *Result)

loop:
	for {
		select {
		case limit <- 1:
			p := RandPerm(n)
			go r.search(p, done)
		case res := <-done:
			resultChan <- res
			<-limit
			// Check if entire solution space traversed
			if r.fs.Full() {
				complete <- true
			}
		case <-stop:
			break loop
		}
	}
}

type Result struct {
	Score float64
	Perm  []uint8
}

type runResult struct {
	Perm   *permutation
	Score  float64
	Opt    bool
	Var    float64
	Center *permutation
	FinalR int
}

func (r *Runner) search(perm *permutation, done chan<- *Result) {
	collect := make(chan *runResult)
	go r.findBestNeighbor(perm, collect)

	// Change what gets sent here
	result := <-collect
	go r.interpret(result, done)
}

func (r *Runner) searchHamming(perm *permutation, dist int, done chan<- *Result) {
	collect := make(chan *runResult)
	go r.sampleHammingRegion(perm, dist, collect)

	// Change what gets sent here
	result := <-collect
	go r.interpret(result, done)
}

// Find best permutation
func (r *Runner) findBestNeighbor(center *permutation, done chan<- *runResult) {
	n := len(center.Seq)
	size := (n * (n - 1) / 2) - 1

	var bestPerm *permutation
	bestScore := math.Inf(1)
	scores := make([]float64, size)

	for i := 0; i < size; i++ {
		neighbor := center.NextNeighbor()
		score := r.Objective(neighbor)
		scores[i] = score

		if score <= bestScore {
			bestScore = score
			bestPerm = neighbor
		}
	}

	isLocalOpt := false
	if centerScore := r.Objective(center); bestScore >= centerScore {
		bestScore = centerScore
		bestPerm = center
		isLocalOpt = true
	}

	// logger.Println(scores)
	vari := variance(scores)

	done <- &runResult{
		Perm:   bestPerm,
		Score:  bestScore,
		Opt:    isLocalOpt,
		Var:    vari,
		Center: center,
		FinalR: 2,
	}
}

// Find best permutation from sampled APPROXIMATE Hamming space
// TODO: predict size of Hamming for max sample size
func (r *Runner) sampleHammingRegion(center *permutation, dist int, done chan<- *runResult) {
	var bestPerm *permutation
	bestScore := math.Inf(1)

	// Determine a reasonable sample size
	n := len(center.Seq)
	size := ((n * (n - 1) / 2) - 1) * dist * dist
	scores := make([]float64, size)

	for i := 0; i < size; i++ {
		neighbor := center.NextHamming(dist)
		score := r.Objective(neighbor)
		scores[i] = score

		if score <= bestScore {
			bestScore = score
			bestPerm = neighbor
		}
	}

	isLocalOpt := false
	if centerScore := r.Objective(center); bestScore >= centerScore {
		bestScore = centerScore
		bestPerm = center
		isLocalOpt = true
	}

	vari := variance(scores)

	done <- &runResult{
		Perm:   bestPerm,
		Score:  bestScore,
		Opt:    isLocalOpt,
		Var:    vari,
		Center: center,
		FinalR: dist,
	}
}

// Use a greedy algorithm search for local mins, but also use stats (variance,
// num time ended up on same path) to determine if to expand the search to a
// greater radius (Hamming distance)
func (r *Runner) interpret(rs *runResult, done chan<- *Result) {
	// Check if already been to the proposed next step
	if r.fs.Test(rs.Perm) {
		// No need to continue further
		logger.Println("Entered previous path")
		done <- nil
		return
	}
	r.fs.Store(rs.Perm)

	// If the solution is optimal, then we're done!
	if rs.Opt {
		logger.Println("Found optimal solution score: ", rs.Score)
		done <- &Result{
			Score: rs.Score,
			Perm:  rs.Perm.Seq,
		}
		return
	}

	// logger.Println("Variance: ", rs.Var)

	// If variance is small look more broadly
	// Otherwise, just follow the best path
	if rs.Var < r.VarCutoff {
		logger.Println("Searching to Hamming dist ", rs.FinalR+1)
		r.searchHamming(rs.Perm, rs.FinalR+1, done)
	} else {
		logger.Println("=> ", rs.Perm)
		r.search(rs.Perm, done)
	}
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
