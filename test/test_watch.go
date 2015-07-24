package main

import (
    "fmt"
)


//共享变量有一个读通道和一个写通道组成
type sharded_var struct {
    reader chan int
    writer chan int
}


//共享变量维护协程
func sharded_var_whachdog(v sharded_var) {
    go func() {
        //初始值
        var value int = 0
        for {
            //监听读写通道，完成服务
            select {
            case value = <-v.writer:
            case v.reader <- value:
            }
        }
    }()
}
func main() {
    //初始化，并开始维护协程
    v := sharded_var{make(chan int), make(chan int)}

    sharded_var_whachdog(v)

    //读取初始值
    fmt.Println(<-v.reader)

    //写入一个值
    v.writer <- 1
    //读取新写入的值
    fmt.Println(<-v.reader)
}