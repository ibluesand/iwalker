package main
import "fmt"


func main() {
    jobs := make(chan int, 5)


    for i :=1 ; i <=6 ; i++ {
        jobs <- i
    }

    fmt.Println(1234)
    fmt.Println(<-jobs)
}
