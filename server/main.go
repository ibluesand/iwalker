package main

import (
	log "code.google.com/p/log4go"
	"github.com/ibluesand/gocomet/server/config"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	Conf *config.Config
)

func main() {

	// init config
	Conf = config.InitConfig()

	// set max routine
	runtime.GOMAXPROCS(Conf.MaxProc)
	// init log
	log.LoadConfiguration(Conf.Log)
	defer log.Close()

	StartServer("7777")

	// init signals, block wait signals
	signalCH := InitSignal()
	HandleSignal(signalCH)
	// exit
	log.Info("comet stop")

}

//错误检查
func checkError(err error, info string) (res bool) {
	if err != nil {
		log.Error(info + "  " + err.Error())
		return false
	}
	return true
}

//启动服务器
func StartServer(port string) {
	service := ":" + port //strconv.Itoa(port);

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")

	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	conns := make(map[string]net.Conn)

	messages := make(chan string, 10)

	//启动服务器广播线程
	go echoHandler(&conns, messages)

	for {
		log.Debug("Listening ...")
		conn, err := l.Accept()
		checkError(err, "Accept")
		log.Debug("Accepting ...")
		conns[conn.RemoteAddr().String()] = conn
		//启动一个新线程
		go Handler(conn, messages)
	}

}

//服务器发送数据的线程
//参数
//      连接字典 conns
//      数据通道 messages
func echoHandler(conns *map[string]net.Conn, messages chan string) {
	for {
		msg := <-messages
		log.Debug(msg)
		for key, value := range *conns {
			log.Debug("connection is connected from ...", key)
			_, err := value.Write([]byte(msg))
			if err != nil {
				log.Debug(err.Error())
				delete(*conns, key)
			}
		}
	}

}

//服务器端接收数据线程
//参数：
//      数据连接 conn
//      通讯通道 messages
func Handler(conn net.Conn, messages chan string) {

	log.Debug("connection is connected from ...", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}
		log.Debug("Rec[", conn.RemoteAddr().String(), "] Say :", string(buf[0:lenght]))
		reciveStr := string(buf[0:lenght])
		messages <- reciveStr
	}
}

// InitSignal register signals handler.
func InitSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	return c
}

// HandleSignal fetch signal from chan then do exit or reload.
func HandleSignal(c chan os.Signal) {
	// Block until a signal is received.
	for {
		s := <-c
		log.Info("comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		// TODO reload
		//return
		default:
			return
		}
	}
}
