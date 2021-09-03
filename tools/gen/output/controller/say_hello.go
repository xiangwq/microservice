package controller

import (
	"context"
	"microservice/tools/gen/output/generate/test"
)

type SayHelloController struct {
}

//检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *SayHelloController) CheckParams(ctx context.Context, r *test.HelloRequest) (err error) {
	return
}

//SayHello函数的实现
func (s *SayHelloController) Run(ctx context.Context, r *test.HelloRequest) (resp *test.HelloResponse, err error) {
	resp = &test.HelloResponse{
		Name: "server reply",
	}
	return resp, nil
}
