package main

import (
    "fmt"
)

func main() {
    q := [3]int{1, 2, 3}
    for _, v := range q {
        fmt.Printf("%d\n", v)
    }
}
