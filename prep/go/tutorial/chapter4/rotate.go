package main

import (
    "fmt"
)

func main() {
    s := []int{0, 1, 2, 3, 4, 5}
    fmt.Println(s)
    rotate_left(s, 2)
    //rotate_right(s, 2)
    fmt.Println(s)
}

// ugh
//func rotate_right(s []int, shift int) {
//    for i, v := range s[shift:] {
//        for j := i + len(s) - shift; j >= 0; j -= shift {
//            s[j], v = v, s[j]
//        }
//    }
//}

func rotate_left(s []int, shift int) {
    length := len(s)
    start := length - shift % length
    // {0, 1, 2, 3, 4, 5}
    // {0, 1, 2, 3, 4, 5}
    // {0, 1, 2, 3, 4, 5}
    // {0, 1, 2, 3, 4, 5}
    // {0, 1, 2, 3, 4, 5}



    // {0, 1, 4, 3, 2, 5}
    // {0, 1, 4, 5, 2, 3}
    // {0, 1, 4, 5, 2, 3}
    // {2, 1, 4, 5, 0, 3}
    // {2, 3, 4, 5, 0, 1}
    var k int
    for i, _ := range s[start:] {
        k = i - shift // maybe fix negative?
        if (k < 0) {
            k = k + length
        }
        s[i], s[k] = s[k], s[i]
    }
}

