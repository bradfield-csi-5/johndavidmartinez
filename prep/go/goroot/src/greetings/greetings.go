package greetings

import "fmt"

// Returns a greeting for named person.
func Hello(name string) string {
    message := fmt.Sprintf("Hi %v. Welcome!", name)
    return message
}
