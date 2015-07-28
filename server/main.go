package main

import (
	log "github.com/cihub/seelog"
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

	requests := make(chan model.Request, 10)

	log.Info("Welcome to chat room ...")

	//启动服务器广播线程
	//go broadcastHandler(messages)

	go redirectHandler(requests)

	var session model.Session
	for {
		conn, err := l.Accept()
		session.Conn = conn
		checkError(err, "Accept")
		//sessions[conn.RemoteAddr().String()] = session
		//启动一个新线程
		go Handler(session, requests)
	}

}

//服务器发送数据的线程
//参数
//      连接字典 conns
//      数据通道 messages
func redirectHandler(messages chan model.Request) {
	for {
		request := <-messages
		var p protocol.Protocol
		p.Type = "broadcast"
		p.Time = time.Now().Unix()

		p.Payload = &request.Payload

		switch  p.Payload.MessageType {

			case "login","all":
				for key, session := range sessions {
					//log.Debug("session.conn[%s], .uid[%s]",session.Conn.RemoteAddr().String(), session.Uid)

					data, err := codec.Eecoder(p)
					if err != nil {
						log.Error(err.Error())
					}
					log.Debugf(string(data))
					_, err = session.Conn.Write(data)
					if err != nil {
						log.Debugf(err.Error())
						delete(sessions, key)
					}
				}

			case "single":

				data, err := codec.Eecoder(p)
				if err != nil {
					log.Error(err.Error())
				}
				log.Debugf(string(data))

				_, err = request.Conn.Write(data)
				if err != nil {
					log.Debugf(err.Error())
					delete(sessions, p.Payload.Uid)
				}

				_, err = sessions[p.Payload.Content.To].Conn.Write(data)
				if err != nil {
					log.Debugf(err.Error())
					delete(sessions, p.Payload.Content.To)
				}




			default:

		}


	}

}

//服务器端接收数据线程
//参数：
//      数据连接 conn
//      通讯通道 messages
func Handler(session model.Session, messages chan model.Request) {

	log.Infof("[%s] join chat room.",session.Conn.RemoteAddr().String())

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
		log.Debugf("[%s] protocol: %s",session.Conn.RemoteAddr().String(), string(buf[0:length]))

		codec.Decoder(buf[0:length], &request)

		var message model.Request
		message.Conn = session.Conn

		switch request.Type {
			case "client":
				switch request.Payload.MessageType {
					case "login":
//						log.Debug(len(sessions), p.Payload.Uid)
//						log.Debug(sessions[session.Conn.RemoteAddr().String()].Conn.RemoteAddr())

						session.Uid = request.Payload.Uid
						sessions[session.Uid] = session

						message.Payload = *request.Payload


					case "all":
						message.Payload = *request.Payload

					case "single":
						message.Payload = *request.Payload

					default:
				}

			default:

		}

		messages <- message
	}
}


