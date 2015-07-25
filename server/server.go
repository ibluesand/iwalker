package main

import (
	log "code.google.com/p/log4go"
	"net"
	"github.com/ibluesand/iwalker/protocol"
	"github.com/ibluesand/iwalker/codec"
	"github.com/ibluesand/iwalker/model"
	"time"
)



func main() {
	StartServer("7777")
}

//错误检查
func checkError(err error, info string) (res bool) {
	if err != nil {
		log.Error(info + "  " + err.Error())
		return false
	}
	return true
}
var sessions = make(map[string]model.Session)

//启动服务器
func StartServer(port string) {
	service := ":" + port //strconv.Itoa(port);

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")

	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	messages := make(chan protocol.Payload, 10)

	log.Info("Welcome to chat room ...")

	//启动服务器广播线程
	go broadcastHandler(messages)

	var session model.Session
	for {
		conn, err := l.Accept()
		session.Conn = conn
		checkError(err, "Accept")
		sessions[conn.RemoteAddr().String()] = session
		//启动一个新线程
		go Handler(session, messages)
	}

}

//服务器发送数据的线程
//参数
//      连接字典 conns
//      数据通道 messages
func broadcastHandler(messages chan protocol.Payload) {
	for {
		msg := <-messages
		var p protocol.Protocol
		p.Type = "broadcast"
		p.Time = time.Now().Unix()

		p.Payload = &msg

		for key, session := range sessions {
			//log.Debug("session.conn[%s], .uid[%s]",session.Conn.RemoteAddr().String(), session.Uid)


			data, err := codec.Eecoder(p)
			if err != nil {
				log.Error(err.Error())
			}
			log.Debug(string(data))
			_, err = session.Conn.Write(data)
			if err != nil {
				log.Debug(err.Error())
				delete(sessions, key)
			}
		}
	}

}

//服务器端接收数据线程
//参数：
//      数据连接 conn
//      通讯通道 messages
func Handler(session model.Session, messages chan protocol.Payload) {

	log.Info("[%s] join chat room.",session.Conn.RemoteAddr().String())

	buf := make([]byte, 1024)

	var request protocol.Protocol
	for {
		length, err := session.Conn.Read(buf)
		if checkError(err, "Connection") == false {
			session.Conn.Close()
			break
		}
		if length > 0 {
			buf[length] = 0
		}
		log.Debug("[%s] protocol: %s",session.Conn.RemoteAddr().String(), string(buf[0:length]))

		codec.Decoder(buf[0:length], &request)

		var message protocol.Payload

		switch request.Type {
			case "client":
				switch request.Payload.MessageType {
					case "login":
//						log.Debug(len(sessions), p.Payload.Uid)
//						log.Debug(sessions[session.Conn.RemoteAddr().String()].Conn.RemoteAddr())
//						sessions[session.Conn.RemoteAddr().String()].Uid = p.Payload.Uid

						session.Uid = request.Payload.Uid
						sessions[session.Conn.RemoteAddr().String()] = session

						message = *request.Payload

					case "message":

						message = *request.Payload

					default:
				}

			default:

		}

		messages <- message
	}
}


