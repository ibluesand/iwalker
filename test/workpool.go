package main

import (
    "fmt"
    "time"
    "math/rand"
    "runtime"
)

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for i :=1 ; i<= runtime.NumCPU(); i++ {
        go worker(i, jobs, results)
    }

    for j := 1; j <= 12; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 12; a++ {
        <-results
    }
}

func worker(id int, jobs <-chan int, results chan<- int) {
    //fmt.Println("I'm a woker ", id)
    for j := range jobs {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))
        t := time.Second * time.Duration(r.Int63n(5))
        fmt.Println("worker", id, "processing job ", j,t)
        time.Sleep(t)
        results <- j * 2
    }
}


