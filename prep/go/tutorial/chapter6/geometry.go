package main

import (
    "math"
    "fmt"
)

type Point struct { X, Y float64 }

type IntPoint struct { X, Y int }

// traditional function
func Distance(p, q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// not allowed
//func Distance(a, b, x, y float64) float64 {
//    return math.Hypot(y - a, x - b)
//}

// method
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}

func (p IntPoint) Distance(q IntPoint) float64 {
    return math.Hypot(float64(q.X)-float64(p.X), float64(q.Y)-float64(p.Y))
}

func main() {
    p := Point{1, 2}
    q := Point{4, 6}

    fmt.Println(Distance(p, q))
    fmt.Println(p.Distance(q))
    ip := IntPoint{1, 2}
    iq := IntPoint{4, 6}

    fmt.Println(ip.Distance(iq))

}
