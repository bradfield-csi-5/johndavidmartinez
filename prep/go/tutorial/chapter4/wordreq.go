package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    wordfreq := make(map[string]int)
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)
    var eof bool
    for {
        eof = scanner.Scan()

        if !eof {
            break
        }

        wordfreq[scanner.Text()]++

        if scanner.Err() != nil {
            fmt.Printf("Scanner Error!\n")
            break
        }
    }
    for word, count := range wordfreq {
        fmt.Printf("%s: %d\n", word, count)
    }
}




