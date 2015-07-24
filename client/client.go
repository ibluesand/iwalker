package main
import (
    "net"
    "fmt"
    "os"
    log "code.google.com/p/log4go"
    "runtime"
    "time"
    "encoding/json"
    "github.com/ibluesand/iwalker/codec"
    "github.com/ibluesand/iwalker/protocol"
)

var uid string

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU())
    conn, err := net.Dial("tcp", "127.0.0.1:7777")
    if err != nil {
        fmt.Println("Error connecting:", err)
        os.Exit(1)
    }
    defer conn.Close()

    login(conn)

    go handleReceive(conn)

    handleSend(conn)
}

func login(conn net.Conn) {

    fmt.Println("please input your id:...")
    _, err := fmt.Scanln(&uid)

    var p protocol.Protocol
    p.ProtocolType = "client"
    p.Time = time.Now().Unix()

    var payload protocol.Payload
    payload.Uid = uid
    payload.MessageType= "login"

    p.Payload = &payload


    login, err := json.Marshal(p)
    if err != nil {
        panic(err.Error())
    }

    conn.Write([]byte(login))

}

func handleReceive(conn net.Conn) {
    buf := make([]byte, 1024)
    var p protocol.Protocol
    for {
        length, err := conn.Read(buf)
        if err != nil {
            conn.Close()
            break
        }
        codec.Decoder(buf[0:length], &p)

        switch p.ProtocolType {
            case "broadcast":
                log.Debug("[%s] login", p.Payload.Content.From)
            default:
        }


    }
}


func handleSend(conn net.Conn) {

    var p protocol.Protocol
    p.ProtocolType = "client"
    p.Time = time.Now().Unix()

    var payload protocol.Payload
    payload.Uid = uid
    payload.MessageType = "message"

    var content protocol.Content
    content.From = uid

    payload.Content = &content

    log.Debug("Welcome!!!")
    for {
        var message string
        _, err := fmt.Scanln(&message)
        if err != nil {
            log.Debug(err.Error())
        }
        content.Message = message
        data, err := codec.Eecoder(p)
        if err != nil {
            log.Debug(err.Error())
        }
        conn.Write([]byte(data))
    }


}
