package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "microservice/example/grpc_example/proto"
)

const address = "localhost:30010"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("dial failed :", address, err.Error())
	}

	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)

	name := "test"

	r, err := c.SayHello(context.TODO(), &pb.HelloRequest{Name: name})

	if err != nil {
		log.Fatalln("request failed: ", err.Error())
	}

	log.Fatalln(r.Name)
}
