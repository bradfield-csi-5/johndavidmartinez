package main

import (
	"fmt"
	"strings"
	"os"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println(strings.Join(os.Args, " "))
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(fmt.Sprintf("Time: %v", elapsed))
}
