package main

import (
    "fmt"
    "example.com/greetings"
)

func main() {
    // Get a greeting and print it
    message := greetings.Hello("John")
    fmt.Println(message)
}
