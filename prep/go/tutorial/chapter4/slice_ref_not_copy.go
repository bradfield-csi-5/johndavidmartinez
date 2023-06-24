package main

import (
    "fmt"
)

// in python a slice is a copy
// in golang it is not
func main() {
    months := []string{0: "Jan", 1: "Feb"}
    month := months[0:1]
    month[0] = "ya"
    fmt.Println(months[0])
}
