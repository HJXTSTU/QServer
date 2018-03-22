package listener

import (
	"net"
)

type AcceptFunc func(conn net.Conn)

type ListenerHandle interface{
	AsyncAccept(onAccept AcceptFunc)
	SyncAccept(onAccept AcceptFunc)
}
type QListener struct {
	listener net.Listener
}


func (this *QListener)accept(onAccept AcceptFunc) {
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			panic(err)
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
