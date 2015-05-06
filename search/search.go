package search

import (
	"github.com/blitzrk/qap-project/matrix"
	"time"
)

type Runner struct {
	NumCPU int
	Cost   matrix.Matrix4D
}

func (r *Runner) Run(stop <-chan int, resultChan chan<- []uint8) {
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
