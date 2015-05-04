package main

import (
	"fmt"
	"github.com/blitzrk/qap-project/dat"
	"io/ioutil"
)

func main() {
	data, err := readDat("bur26a.dat")
	if err != nil {
	  fmt.Println(err)
	  return
	}
	
	n := data[0][0][0]
	times := data[1]
	freqs := data[2]
	totalF := sum(freqs)

	fmt.Println(n)
	fmt.Println()
	fmt.Println(times)
	fmt.Println()
	fmt.Println(freqs)
	fmt.Println()
	fmt.Println(totalF)
	// 	fmt.Println()
	// 	fmt.Println(apply(freqs, func(x float64) float64 { return n*n*x/totalF }))
}

func readDat(fname string) ([][][]float64, error) {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
  
  return dat.Read(f), nil
}
