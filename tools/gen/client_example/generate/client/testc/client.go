package testc

import (
	"context"
	"fmt"

	"microservice/tools/gen/client_example/generate/test"

	"microservice/errno"
	"microservice/meta"
	"microservice/rpc"
)

type TestClient struct {
	serviceName string
	client      *rpc.MicroserviceClient
}

func NewTestClient(serviceName string, opts ...rpc.RpcOptionFunc) *TestClient {
	c := &TestClient{
		serviceName: serviceName,
	}
	c.client = rpc.NewMicroserviceClient(serviceName, opts...)
	return c
}

func (s *TestClient) SayHello(ctx context.Context, r *test.HelloRequest) (resp *test.HelloResponse, err error) {
	/*
		middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
		mkResp, err := middlewareFunc(ctx, r)
		if err != nil {
			return nil, err
		}
	*/
	mkResp, err := s.client.Call(ctx, "SayHello", r, mwClientSayHello)
	if err != nil {
		return nil, err
	}
	resp, ok := mkResp.(*test.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *test.HelloResponse")
		return nil, err
	}

	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	/*
		conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
		if err != nil {
			logs.Error(ctx, "did not connect: %v", err)
			return nil, err
		}*/
	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*test.HelloRequest)
	client := test.NewHelloServiceClient(rpcMeta.Conn)

	return client.SayHello(ctx, req)
}
