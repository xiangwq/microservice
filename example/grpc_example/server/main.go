package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "microservice/example/grpc_example/proto"
	"net"
)

const port = ":30010"

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Name: "hello:" + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("failed to listen", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	s.Serve(lis)
}
