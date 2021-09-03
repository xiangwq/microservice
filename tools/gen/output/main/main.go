package main

import (
	"log"
	"microservice/server"

	"microservice/tools/gen/output/generate/test"
	"microservice/tools/gen/output/router"
)

var routerServer = &router.RouterServer{}

func main() {
	err := server.Init("test")
	if err != nil {
		log.Fatal("init service failed, err: $v", err)
		return
	}

	test.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
