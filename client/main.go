package main
import (
    "net"
    "fmt"
    "os"
    log "github.com/cihub/seelog"
    "runtime"
    "time"
    "encoding/json"
    "github.com/ibluesand/iwalker/codec"
    "github.com/ibluesand/iwalker/protocol"
    "strings"
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

    var request protocol.Protocol
    request.Type = "client"
    request.Time = time.Now().Unix()

    var payload protocol.Payload
    payload.Uid = uid
    payload.MessageType= "login"

    request.Payload = &payload


    login, err := json.Marshal(request)
    if err != nil {
        log.Error(err.Error())
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

        switch p.Type {
            case "broadcast":
                switch p.Payload.MessageType {
                    case "login":
                        log.Debugf("[%s] login", p.Payload.Uid)
                    case "logout":
                        log.Debugf("[%s] logout", p.Payload.Uid)
                    case "all":
                        log.Debugf("[%s] says: [%s]", p.Payload.Uid, p.Payload.Content.Message)
                    case "single":
                        log.Debugf("[%s]->[%s] says: [%s]", p.Payload.Uid, p.Payload.Content.To, p.Payload.Content.Message)
                    default:
                        log.Debugf(string(buf[0:length]))
                }

            default:
        }

    }
}


func handleSend(conn net.Conn) {

    var p protocol.Protocol
    p.Type = "client"
    p.Time = time.Now().Unix()

    var payload protocol.Payload
    payload.Uid = uid

    var content protocol.Content
    content.From = uid

    p.Payload = &payload
    payload.Content = &content

    log.Debug("Welcome!!!")
    for {
        var message string
        _, err := fmt.Scanln(&message)
        if err != nil {
            log.Debug(err.Error())
        }

        index := strings.Index(message, ":")
        log.Debug(index)
        if strings.HasPrefix(message, "@") && index > 1{
            payload.MessageType = "single"
            content.Message = string([]rune(message)[index+1:])
            content.To=string([]rune(message)[1:index])
        } else {
            payload.MessageType = "all"
            content.Message = message
        }

        data, err := codec.Eecoder(p)
        if err != nil {
            log.Debug(err.Error())
        }
        conn.Write([]byte(data))
    }


}
