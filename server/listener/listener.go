package listener

import (
	"net"
)

const(
	DEFAULT_LIMIT = 1024
)

type AcceptFunc func(conn net.Conn)

var connLimitChanel chan struct{}

func init() {
	connLimitChanel = make(chan struct{},DEFAULT_LIMIT)
}

type ListenerHandle interface{
	AsyncAccept(onAccept AcceptFunc)
	SyncAccept(onAccept AcceptFunc)
	Close()
	ReleaseConn()
}

type QListener struct {
	listener net.Listener
}

func (this *QListener)Close(){
	this.listener.Close()
}

func (this *QListener)ReleaseConn(){
	<-connLimitChanel
}


func (this *QListener)accept(onAccept AcceptFunc) {
	for {
		connLimitChanel<- struct{}{}
		conn, err := this.listener.Accept()
		if err != nil {
			break
		}
		go onAccept(conn)
	}
}

func (this *QListener) AsyncAccept(onAccept AcceptFunc) {
	go this.accept(onAccept)
}

func (this *QListener) SyncAccept(onAccept AcceptFunc)  {
	this.accept(onAccept)
}

func NewListener(address string) ListenerHandle {
	listener := QListener{}
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	listener.listener = l
	return &listener
}
