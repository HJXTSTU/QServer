package main

import (
	"fmt"
	"projects/socket_component/server"
	"time"
	"projects/socket_component/server/connection"
	"projects/socket_component/util"
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
	tp := TopProcesser{}
	s := server.NewQServer(":8888")
	s.SetProcesser(&tp)
	s.Listen()
	fmt.Println("Server Boot")
	for{
		time.Sleep(time.Second)
	}

}