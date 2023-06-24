package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + fmt.Sprintf("%v: %s", i, os.Args[i])
		sep = " "
	}
	fmt.Println(s)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(fmt.Sprintf("Time: %v", elapsed))
}
