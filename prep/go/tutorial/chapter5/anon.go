package main

import "fmt"

func squares() func() int {
    var x int
    return func() int {
        x++
        return x * x;
    }
}

func main() {
    f := squares()
    for i := 0; i < 4; i++ {
        fmt.Println(f())
    }
}
