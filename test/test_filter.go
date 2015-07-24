package main

import (
//    "fmt"
    "fmt"
)

func generator(ch chan<- int) {
    for i := 2; ; i++ {
        fmt.Println("i",i)
        ch<- i // Send 'i' to channel 'ch'.
    }
}

// Copy the values from channel 'in' to channel 'out',
//removing those divisible by 'prime'.
func filter(in <-chan int, out chan<- int, prime int) {
    for {
        i := <-in // Receive value from 'in'.
        if i%prime != 0 {
            out <- i // Send 'i' to 'out'.
        }
    }
}

func main() {
    ch := make(chan int) // Create a new channel.

    go generator(ch)      // Launch Generate goroutine.

    for i := 0; i < 10; i++ {
        prime := <-ch
        print(prime, "\n")
        ch1 := make(chan int)
        go filter(ch, ch1, prime)
        ch = ch1
    }
}

