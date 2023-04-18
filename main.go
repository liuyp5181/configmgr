package main

import (
	"github.com/liuyp5181/base/config"
	"github.com/liuyp5181/base/database"
	"github.com/liuyp5181/base/service"
	pb "github.com/liuyp5181/configmgr/api"
	"github.com/liuyp5181/configmgr/handler"
)

func init() {
	config.ServiceName = pb.Greeter_ServiceDesc.ServiceName
}

func main() {

	err := database.Connect("test")
	if err != nil {
		panic(err)
	}

	err = service.InitClients()
	if err != nil {
		panic(err)
	}

	s := service.NewServer()
	pb.RegisterGreeterServer(s, &handler.ServiceImpl{})
	s.Serve()
}
