package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"microservice/logs"
	"microservice/rpc"
	pb "microservice/tools/gen/output/generate"
)

type HelloClient struct {
	serviceName string
}

func NewHelloClient(serviceName string) *HelloClient {
	return &HelloClient{
		serviceName: serviceName,
	}
}

func (h *HelloClient) SayHelloV1(ctx context.Context, in *pb.HelloRequest, opts ...grpc.CallOption) (*pb.HelloResponse, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(context.Background(), "did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)
	r, err := c.SayHello(ctx, in, opts...)
	if err != nil {
		logs.Error(ctx, "could not greet: %v", err)
		return nil, err
	}
	return r, err
}

func (h *HelloClient) SayHello(ctx context.Context, in *pb.HelloRequest, opts ...grpc.CallOption) (*pb.HelloResponse, error) {

	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, in)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*pb.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
		return nil, err
	}

	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}
	req := request.(*pb.HelloRequest)
	defer conn.Close()
	client := pb.NewHelloServiceClient(conn)
	return client.SayHello(ctx, req)
}
