package main

import (
    "fmt"
)

func main() {
    stuff := []string{"pizza", "bagel", "pizza", "pizza", "pizza"}
    fmt.Println(stuff)
    dadup := dedup(stuff)
    fmt.Println(stuff)
    fmt.Println(dadup)
}

func dedup(strings []string) []string {
    rmc := 0
    for i := 1; i < len(strings); i++ {
        if same(strings[i], strings[i - 1]) {
            strings[i - 1] = "" 
            rmc++
        }
    }
    for k, i := 0, 0; k < len(strings); {
        strings[i], strings[k] = strings[k], strings[i]
        if strings[i] != "" {
            i++
        }
        k++
    }
    return strings[:len(strings)-rmc]
}

func same(x string, y string) bool {
    if (len(x) != len(y)) {
        return false
    }
    i := len(x) - 1
    for {
        if y[i] != x[i] {
            return false
        }
        if i == 0 {
            break
        }
        i--
    }
    return true
}
