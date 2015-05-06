package search

import (
	"fmt"
	"time"
)

func Run(n, CPUs int, stop <-chan int) {
	limit := make(chan int, CPUs)
	done := make(chan int)

	for {
		select {
		case limit <- 1:
			go search(n, done)
		case <-done:
			<-limit
		case <-stop:
			break
		}
	}
}

func search(n int, done chan<- int) {
	p := RandPerm(n)
	fmt.Println()
	fmt.Println(p)
	time.Sleep(1 * time.Minute)
	done <- 1
}
