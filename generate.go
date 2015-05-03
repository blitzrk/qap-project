package main

import (
	"fmt"
	"github.com/blitzrk/generate/dat"
	"io/ioutil"
)

func main() {
	f, err := ioutil.ReadFile("bur26a.dat")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := dat.Read(f)
	n := data[0][0][0]
	times := data[1]
	freqs := data[2]

	fmt.Println(n)
	fmt.Println()
	fmt.Println(times)
	fmt.Println()
	fmt.Println(freqs)
}
