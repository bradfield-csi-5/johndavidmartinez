package main

import (
    "flag"
    "fmt"
    "os"
    "bufio"
    "strings"
    "crypto/sha256"
    "crypto/sha512"
)

var hashFunc = flag.String("c", "sha256", "hash function")

func main() {
    flag.Parse()
    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading from stdin")
    }
    trimmed := strings.TrimSuffix(text, "\n")

    switch *hashFunc {
    case "sha256":
        hashed := sha256.Sum256([]byte(trimmed))
        fmt.Printf("%x", hashed)
    case "sha384":
        hashed := sha512.Sum384([]byte(trimmed))
        fmt.Printf("%x", hashed)
    case "sha512":
        hashed := sha512.Sum512([]byte(trimmed))
        fmt.Printf("%x", hashed)
    default:
        fmt.Println("Unsupported hash function")
        os.Exit(1)
    }
}

