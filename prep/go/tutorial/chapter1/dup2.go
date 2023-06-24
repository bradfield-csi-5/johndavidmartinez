package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	file_names := make(map[string]map[string]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLinesStdin(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, file_names, arg)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			files_for_line := ""
			for file_for_line, _ := range file_names[line] {
				files_for_line += file_for_line + " "
			}
			fmt.Printf("file: %v | %d\t%s\n", files_for_line, n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int, file_names map[string]map[string]string, file_name string) {
	input := bufio.NewScanner(f)
        for input.Scan() {
		txt := input.Text()
		//file_names[txt] = make(map[string]string)
		if file_names[txt] == nil {
			file_names[txt] = make(map[string]string)
		}
		file_names[txt][file_name] = ""
		counts[txt]++
        }
}

func countLinesStdin(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
        for input.Scan() {
		txt := input.Text()
		counts[txt]++
        }
}
