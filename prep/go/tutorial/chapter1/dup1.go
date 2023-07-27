package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	INDEX_CMD  = "index"
	SEARCH_CMD = "search"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	// eats errors
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
