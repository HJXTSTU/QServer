package main

import (
	"fmt"
	"projects/socket_component/server/connection"
	"wwt/util"
	"projects/socket_component/server"
	"projects/socket_component/ctrl"
	"os"
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
	stream.InsertLen()
	token.SendAsync(stream.Bytes(), func(handler connection.TokenHandler, bytes []byte, i int, e error) {
		fmt.Println("Send Successful")
	})
}



func main() {

	s := server.NewQServer(":8888")
	tp := TopProcesser{}
	s.SetProcesser(&tp)
	s.AsyncListen()

	goroutine_group := ctrl.GlobalWaitGroup()
	exit_chan := ctrl.GlobalExitChan()
	<-exit_chan
	s.Close()
	goroutine_group.Wait()
	os.Exit(0)
}