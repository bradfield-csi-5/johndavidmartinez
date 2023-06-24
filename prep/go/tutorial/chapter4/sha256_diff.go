package main

import (
    "fmt"
    "crypto/sha256"
)

func main() {
    c1 := sha256.Sum256([]byte("x"))
    c2 := sha256.Sum256([]byte("X"))
    fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
    fmt.Println("---------")
    fmt.Println("diff: %v", sha256Diff([]byte("x"), []byte("X")))
}

func sha256Diff(x []byte, y []byte) int {
    xs := sha256.Sum256(x)
    ys := sha256.Sum256(y)
    diffcount := 0
    for i, _ := range xs {
        if xs[i] != ys[i] {
            diffcount++
        }
    }
    return diffcount
}
