package main

import (
	"fmt"
	"socket_component/server/connection"
	"socket_component/util"
	"socket_component/server"
)

type TopProcesser struct{

}

func (this *TopProcesser)Processe(token connection.TokenHandler,length int,bytes []byte){
	stream := util.NewStreamBuffer()
	stream.Append(bytes)
	i := stream.ReadLine()
	fmt.Println(i)
	stream.Renew()
	stream.WriteLine("LO")
	sd:=util.NewStreamBuffer()
	sd.WriteInt(stream.Len())
	sd.WriteNBytes(stream.Bytes(),stream.Len())
	token.SendAsync(sd.Bytes(), func(handler connection.TokenHandler, bytes []byte, i int, e error) {
		if e == nil{
			fmt.Printf("Send sucessful.\n")
		}
	})
}



func main() {

	s := server.NewQServer(":8888")
	tp := TopProcesser{}
	s.SetProcesser(&tp)
	s.SyncListen()

	fmt.Println("test sync")

}