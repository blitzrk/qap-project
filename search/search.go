package search

import (
	"fmt"
	"github.com/blitzrk/qap-project/matrix"
	"time"
)

type Runner struct {
	NumCPU int
	Cost   matrix.Matrix4D
	n      int
}

func (r *Runner) Run(stop <-chan int) {
	limit := make(chan int, r.NumCPU)
	done := make(chan int)

	for {
		select {
		case limit <- 1:
			go r.search(done)
		case <-done:
			<-limit
		case <-stop:
			break
		}
	}
}

func (r *Runner) search(done chan<- int) {
	r.n = len(r.Cost)
	p := RandPerm(r.n)

	fmt.Println()
	fmt.Println(p)
	time.Sleep(1 * time.Minute)

	done <- 1
}
