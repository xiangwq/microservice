package router

import (
	"context"
	"microservice/meta"
	"microservice/server"

	"microservice/tools/gen/output/controller"
	"microservice/tools/gen/output/generate/test"
)

type RouterServer struct{}

func (s *RouterServer) SayHello(ctx context.Context, r *test.HelloRequest) (resp *test.HelloResponse, err error) {
	ctx = meta.InitServerMeta(ctx, "test", "SayHello")
	mwFunc := server.BuildServerMiddleware(mwSayHello)
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	resp = mwResp.(*test.HelloResponse)
	return resp, err
}

func mwSayHello(ctx context.Context, req interface{}) (resp interface{}, err error) {
	r := req.(*test.HelloRequest)
	ctrl := &controller.SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}
