package model

import (
    "net"
    "github.com/ibluesand/iwalker/protocol"
)

type Session struct {
    Conn net.Conn
    Uid string
}

type Request struct {
    Conn net.Conn
    Payload protocol.Payload
}
