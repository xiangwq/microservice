package server

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"microservice/logs"
	"microservice/middleware"
	"microservice/registry"
	_ "microservice/registry/etcd"
	"microservice/util"
	"net"
	"net/http"
)

type MicroserviceServer struct {
	*grpc.Server
	limiter        *rate.Limiter
	register       registry.Registry
	userMiddleware []middleware.Middleware
}

var microserviceServer = &MicroserviceServer{
	Server: grpc.NewServer(),
}

func Use(m ...middleware.Middleware) {
	microserviceServer.userMiddleware = append(microserviceServer.userMiddleware, m...)
}

func Init(serviceName string) error {
	err := InitConfig(serviceName)
	if err != nil {
		return err
	}
	initLogger()

	//初始化注册中心
	err = initRegister(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init register failed, err:%v", err)
		return err
	}

	err = initTrace(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init tracing failed, err:%v", err)
	}

	//初始化限流器
	if microserviceConf.Limit.SwitchOn {
		microserviceServer.limiter = rate.NewLimiter(rate.Limit(microserviceConf.Limit.QPSLimit), microserviceConf.Limit.QPSLimit)
	}

	return nil
}

func initLogger() (err error) {
	filename := fmt.Sprintf("%s/%s.log", microserviceConf.Log.Dir, microserviceConf.ServiceName)
	outputer, err := logs.NewFileOutputer(filename)
	if err != nil {
		return
	}

	level := logs.GetLogLevel(microserviceConf.Log.Level)
	logs.InitLogger(level, microserviceConf.Log.ChanSize, microserviceConf.ServiceName)
	logs.AddOutputer(outputer)

	if microserviceConf.Log.ConsoleLog {
		logs.AddOutputer(logs.NewConsoleOutputer())
	}
	return
}

func initTrace(serviceName string) (err error) {

	if !microserviceConf.Trace.SwitchOn {
		return
	}

	return middleware.InitTrace(serviceName, microserviceConf.Trace.ReportAddr,
		microserviceConf.Trace.SampleType, microserviceConf.Trace.SampleRate)
}

func initRegister(serviceName string) (err error) {

	if !microserviceConf.Register.SwitchOn {
		return
	}

	ctx := context.TODO()
	fmt.Printf(microserviceConf.Register.RegisterName)
	registryInst, err := registry.InitRegistry(ctx,
		microserviceConf.Register.RegisterName,
		registry.WithAddrs([]string{microserviceConf.Register.RegisterAddr}),
		registry.WithTimeout(microserviceConf.Register.Timeout),
		registry.WithRegistryPath(microserviceConf.Register.RegisterPath),
		registry.WithHeartBeat(microserviceConf.Register.HeartBeat),
	)
	if err != nil {
		logs.Error(ctx, "init registry failed, err:%v", err)
		return
	}

	microserviceServer.register = registryInst
	service := &registry.Service{
		Name: serviceName,
	}

	ip, err := util.GetLocalIP()
	if err != nil {
		return
	}
	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   ip,
		Port: uint(microserviceConf.Port),
	},
	)

	registryInst.Register(context.TODO(), service)
	return
}

func BuildServerMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	fmt.Println("test middle")
	var mids []middleware.Middleware
	mids = append(mids, middleware.AccessLogMiddleware)

	if microserviceConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	if microserviceConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewRateLimitMiddleware(microserviceServer.limiter))
	}

	if microserviceConf.Trace.SwitchOn {
		fmt.Println("trace middle is append")
		mids = append(mids, middleware.TraceServerMiddleware)
	}

	m := middleware.Chain(middleware.PrepareMiddleware, mids...)

	return m(handle)
}

func Run() {
	if GetConf().Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", GetConf().Prometheus.Port)
			http.ListenAndServe(addr, nil)
		}()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GetServerPort()))
	if err != nil {
		log.Fatalln("failed to listen", err.Error())
	}
	microserviceServer.Serve(lis)
}

func GRPCServer() *grpc.Server {
	return microserviceServer.Server
}
