package main
import (
    "fmt"
    "time"
)
func main() {
    go sayHello()
    time.Sleep(time.Second)
}

func sayHello() {
    fmt.Println("say hello")
    go sayHandleHello()
}

func sayHandleHello() {
    fmt.Println("sayHandleHello")
}