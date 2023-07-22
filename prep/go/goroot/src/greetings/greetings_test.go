package greetings

import "testing"

func TestGreeting(t *testing.T) {
    want := "Hi John. Welcome!"
    if got := Hello("John"); got != want {
        t.Error("Not working")
    }
}
