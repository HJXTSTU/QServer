package listener

import (
	"net"
)


type AcceptFunc func(conn net.Conn)

type ListenerHandle interface{
	AcceptAsync(onAccept AcceptFunc)
}
type QListener struct {
	listener net.Listener
}

func (this *QListener) AcceptAsync(onAccept AcceptFunc) {
	go func(onAccept AcceptFunc) {
		for {
			conn, err := this.listener.Accept()
			if err != nil {
				panic(err)
			}
			go onAccept(conn)
		}
	}(onAccept)
}

func NewListener(address string) ListenerHandle {
	listener := QListener{}
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	listener.listener = l;
	return &listener
}
