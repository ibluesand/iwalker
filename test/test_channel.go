package main

import (
    "fmt"
    "math/rand"
)

func rand_generator() chan int {
    out := make(chan int)

    go func() {
       for {
           out <- rand.Int()
       }
    }()

    return out;
}

func rand_generate_muti() chan int {
    rand_generator_1 := rand_generator()
    rand_generator_2 := rand_generator()

    //创建通道
    out := make(chan int)

    //创建协程
    go func() {
        for {
            //读取生成器1中的数据，整合
            out <- <-rand_generator_1
        }
    }()
    go func() {
        for {
            //读取生成器2中的数据，整合
            out <- <-rand_generator_2
        }
    }()
    return out
}

func main() {


    // 生成随机数作为一个服务
    rand_service_handler :=rand_generator()
    // 从服务中读取随机数并打印
    fmt.Printf("%d\n",<-rand_service_handler)

    fmt.Printf("%d\n",<-rand_service_handler)


}

