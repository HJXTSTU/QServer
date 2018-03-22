package server

import (
	"projects/socket_component/server/listener"
	"projects/socket_component/server/connection"
	"net"
)

type QServerHandle interface {
	Listen()

	SetProcesser(QProcesser)
}

type QProcesser interface {
	Processe(connection.TokenHandler, int, []byte)
}

type QWriter interface {
	Send([]byte)
}

type QServer struct {
	listener  listener.ListenerHandle
	tokens    connection.TokenPoolHandler
	processer QProcesser
}

func (this *QServer) Listen() {
	this.listener.AcceptAsync(this.onAccept)
}

func (this *QServer) onAccept(conn net.Conn) {
	token := connection.NewQToken(conn, this.onRead, this.onClose)
	this.tokens.AddToken(token)
	token.ReadAsync()
	//this.listener.AcceptAsync(this.onAccept)
}

func (this *QServer) onRead(handle connection.TokenHandler, n int, bytes []byte) {
	this.processer.Processe(handle, n, bytes)
}

func (this *QServer) SetProcesser(p QProcesser) {
	this.processer = p
}

func (this *QServer) onClose(handle connection.TokenHandler) {
	//TODO::关闭TOKEN
	handle.Close()
	this.tokens.DeleteToken(handle)
}

func NewQServer(address string) QServerHandle {
	qserver := new(QServer)
	qserver.listener = listener.NewListener(address)
	qserver.tokens = connection.NewTokenPool()
	return qserver
}
