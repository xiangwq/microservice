package rpc

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"microservice/loadbalance"
	"microservice/logs"
	"microservice/meta"
	"microservice/middleware"
	"microservice/registry"
)

var initRegistryOnce sync.Once
var globalRegister registry.Registry

type MicroserviceClient struct {
	opts     *RpcOptions
	register registry.Registry
	limiter  *rate.Limiter
	balance  loadbalance.LoadBalance
}

func NewMicroserviceClient(serviceName string, optfunc ...RpcOptionFunc) *MicroserviceClient {
	client := &MicroserviceClient{
		opts: &RpcOptions{
			ConnTimeout:       DefaultConnTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			ReadTimeout:       DefaultReadTimeout,
			ServiceName:       serviceName,
			RegisterName:      "etcd",
			RegisterAddr:      "127.0.0.1:2379",
			RegisterPath:      "/microservice/service/",
			TraceReportAddr:   "http://127.0.0.1:9412/api/v1/spans",
			TraceSampleType:   "const",
			TraceSampleRate:   1,
			ClientServiceName: "default",
		},
		balance: loadbalance.NewRandomBalance(),
	}

	for _, opt := range optfunc {
		opt(client.opts)
	}

	initRegistryOnce.Do(func() {
		ctx := context.TODO()
		var err error
		globalRegister, err = registry.InitRegistry(ctx,
			client.opts.RegisterName,
			registry.WithAddrs([]string{client.opts.RegisterAddr}),
			registry.WithTimeout(time.Second),
			registry.WithRegistryPath(client.opts.RegisterPath),
			registry.WithHeartBeat(10),
		)
		if err != nil {
			logs.Error(ctx, "init registry failed, err:%v", err)
			return
		}
	})

	if client.opts.MaxLimitQps > 0 {
		client.limiter = rate.NewLimiter(rate.Limit(client.opts.MaxLimitQps),
			client.opts.MaxLimitQps)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf("0.0.0.0:%d", 8082)
		http.ListenAndServe(addr, nil)
	}()

	err := middleware.InitTrace(client.opts.ClientServiceName, client.opts.TraceReportAddr, client.opts.TraceSampleType,
		client.opts.TraceSampleRate)

	if err != nil {
		fmt.Println("init trace failed", err.Error())
	}
	client.register = globalRegister
	return client
}

func (k *MicroserviceClient) getCaller(ctx context.Context) string {

	serverMeta := meta.GetServerMeta(ctx)
	if serverMeta == nil {
		return ""
	}
	return serverMeta.ServiceName
}

func (k *MicroserviceClient) buildMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {

	var mids []middleware.Middleware
	mids = append(mids, middleware.PrepareMiddleware)
	mids = append(mids, middleware.RpcLogMiddleware)
	mids = append(mids, middleware.TraceRpcMiddleware)
	mids = append(mids, middleware.PrometheusRpcMiddleware)
	if k.limiter != nil {
		mids = append(mids, middleware.NewRateLimitMiddleware(k.limiter))
	}

	mids = append(mids, middleware.HystrixMiddleware)
	mids = append(mids, middleware.NewDiscoveryMiddleware(k.register))

	mids = append(mids, middleware.NewLoadBalanceMiddleware(k.balance))
	mids = append(mids, middleware.ShortConnectMiddleware)

	m := middleware.Chain(mids[0], mids[1:]...)
	return m(handle)

}

func (k *MicroserviceClient) Call(ctx context.Context, method string, r interface{}, handle middleware.MiddlewareFunc) (resp interface{}, err error) {

	//构建中间件
	caller := k.getCaller(ctx)
	ctx = meta.InitRpcMeta(ctx, k.opts.ServiceName, method, caller)
	fmt.Println("build middle")
	middlewareFunc := k.buildMiddleware(handle)
	resp, err = middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp, err
}
