package main

import (
	log "github.com/cihub/seelog"
	"net"
	"github.com/ibluesand/iwalker/protocol"
	"github.com/ibluesand/iwalker/codec"
	"github.com/ibluesand/iwalker/model"
	"time"
	"github.com/ibluesand/iwalker/common"
)



func main() {
	StartServer("7777")
}


var sessions = make(map[string]model.Session)

//启动服务器
func StartServer(port string) {
	service := ":" + port //strconv.Itoa(port);

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	common.CheckError(err, "ResolveTCPAddr")

	l, err := net.ListenTCP("tcp", tcpAddr)
	common.CheckError(err, "ListenTCP")

	requests := make(chan model.Request, 10)

	log.Info("Welcome to chat room ...")

	//启动服务器广播线程
	//go broadcastHandler(messages)

	go redirectHandler(requests)

	var session model.Session
	for {
		conn, err := l.Accept()
		session.Conn = conn
		common.CheckError(err, "Accept")
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


		data, err := codec.Eecoder(p)
		if err != nil {
			log.Error(err.Error())
		}
		log.Debugf("[debug] [%s]",string(data))


		switch  p.Payload.MessageType {

			case "login","logout","all":
				for key, session := range sessions {
					//log.Debug("session.conn[%s], .uid[%s]",session.Conn.RemoteAddr().String(), session.Uid)

					_, err = session.Conn.Write(data)
					if err != nil {
						log.Debugf(err.Error())
						delete(sessions, key)
					}
				}

			case "single":

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

		var message model.Request
		message.Conn = session.Conn

		//
		if common.CheckError(err, "Connection error") == false {

			log.Debug("connect error")
			var payload protocol.Payload
			payload.Uid = session.Uid
			payload.MessageType = "logout"
			message.Payload = payload

			session.Conn.Close()
			delete(sessions, session.Uid)

			log.Debug("message logout->::", message.Payload.MessageType)
			messages <- message

			break

		} else {
			log.Debugf("[%s] protocol: %s", session.Conn.RemoteAddr().String(), string(buf[0:length]))

			codec.Decoder(buf[0:length], &request)

			switch request.Type {
				case "client":
				switch request.Payload.MessageType {
					case "login":
					session.Uid = request.Payload.Uid
					sessions[session.Uid] = session

					case "all":

					case "single":

					default:
				}
				message.Payload = *request.Payload

				default:

			}

			log.Debug("message->::", message.Payload.MessageType)
			messages <- message
		}


	}
}


