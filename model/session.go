package model

import (
    "net"
)

type Session struct {
    Conn net.Conn
    Uid string
}
