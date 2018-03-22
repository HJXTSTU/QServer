package connection

import (
	"net"
	"socket_component/util"
	"sync"
)

const (
	BUFFER_SIZE = 1024
)

type ReadCallback func(TokenHandler, int, []byte)
type CloseCallback func(TokenHandler)
type SendCallback func(TokenHandler, []byte, int, error)

type TokenPoolHandler interface {
	AddToken(token TokenHandler)
	DeleteToken(token TokenHandler)
	Len() int
}

type TokenPool struct {
	tokens map[TokenHandler]TokenHandler
	mu     sync.Mutex
}

func (this *TokenPool) AddToken(token TokenHandler) {
	this.mu.Lock()
	this.tokens[token] = token
	this.mu.Unlock()
}

func (this *TokenPool) DeleteToken(token TokenHandler) {
	this.mu.Lock()
	delete(this.tokens, token)
	this.mu.Unlock()
}

func (this *TokenPool) Len() int {
	this.mu.Lock()
	l := len(this.tokens)
	this.mu.Unlock()
	return l
}

func NewTokenPool() TokenPoolHandler {
	return &TokenPool{make(map[TokenHandler]TokenHandler), sync.Mutex{}}
}

type TokenHandler interface {
	ReadAsync()

	Close()

	OnRead(TokenHandler, int, []byte)

	OnClose(TokenHandler)

	RemoteAddr() net.Addr

	SendAsync(b []byte, callback SendCallback)
}

type QToken struct {
	conn     net.Conn
	onRead   ReadCallback
	onClose  CloseCallback
	r_stream util.StreamBuffer
}

func (this *QToken) SendAsync(b []byte, callback SendCallback) {
	w_stream := util.NewStreamBuffer()
	w_stream.WriteInt(len(b))
	w_stream.Append(b)
	n, err := this.conn.Write(w_stream.Bytes())
	if callback != nil {
		callback(this, b, n, err)
	}
}

func (this *QToken) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *QToken) ReadAsync() {
	go func(handle TokenHandler) {
		for {
			buf := make([]byte, BUFFER_SIZE)
			n, err := this.conn.Read(buf)
			if n <= 0 || err != nil {
				break;
			}
			handle.OnRead(handle, n, buf[:n])
		}
		// TODO::onClose
		handle.OnClose(handle)
	}(this)
}

func (this *QToken) OnRead(handle TokenHandler, n int, bytes []byte) {
	this.r_stream.Append(bytes)
	lstream := this.r_stream.Len()
	if lstream < 4 {
		return
	}
	length := this.r_stream.ReadInt()
	//	数据包已经完整
	if length <= lstream {
		data := this.r_stream.ReadNBytes(length)
		go this.onRead(this, length, data)
	} else {
		this.r_stream.Undo()
	}

}

func (this *QToken) OnClose(handle TokenHandler) {
	go this.onClose(this)
}

func (this *QToken) Close() {
	this.conn.Close()
}

func NewQToken(conn net.Conn, onRead ReadCallback, onClose CloseCallback) *QToken {
	token := QToken{conn, onRead, onClose, util.NewStreamBuffer()}
	return &token
}
