package config

import (
    "runtime"
)

type Config struct {
    // base section
    PidFile   string
    Dir       string
    Log       string
    MaxProc   int

    // proto section
    TCPBind          []string
    TCPSndbuf        int
    TCPRcvbuf        int
    TCPKeepalive     bool

    // channel
    SndbufSize            int
    RcvbufSize            int

    BufioInstance int
    BufioNum int

}

func InitConfig() *Config {
    return &Config{
        // base section
        PidFile:   "/tmp/gopush-cluster-comet.pid",
        Dir:       "./",
        Log:       "/Users/bluesand/gospace/bin/comet_log.xml",
        MaxProc:   runtime.NumCPU(),

        // proto section
        TCPBind:          []string{"localhost:7777"},
        TCPSndbuf:        1024,
        TCPRcvbuf:        1024,
        TCPKeepalive:     false,

        // channel
        SndbufSize:              2048,
        RcvbufSize:              256,

        BufioInstance:           runtime.NumCPU(),
        BufioNum:                128,
    }
}



