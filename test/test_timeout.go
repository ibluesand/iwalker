package main

import (
    "time"
)

func never_leak(ch chan int) {
    //初始化timeout，缓冲为1
    timeout := make(chan bool, 1)
    //启动timeout协程，由于缓存为1，不可能泄露
    go func() {
        time.Sleep(1 * time.Second)
        timeout <- true
    }()
    //监听通道，由于设有超时，不可能泄露
    select {
    case <-ch:
    // a read from ch hasoccurred
    case <-timeout:
    // the read from ch has timedout
    }
}

func main() {
}