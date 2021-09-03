package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"io"
	"microservice/logs"
	"os"
	"time"
)

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(service string) (opentracing.Tracer, io.Closer) {

	transport, err := zipkin.NewHTTPTransport(
		"http://127.0.0.1:9412/api/v1/spans",
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init zipkin: %v\n", err))
	}

	fmt.Printf("transport:%v\n", transport)
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)
	fmt.Printf("r=%v\n", r)
	tracer, closer, err := cfg.New(service,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	tracer, err := InitTrace("6", "http://127.0.0.1:9412/api/v1/spans", "const", 1)

	if err != nil {
		fmt.Println(err.Error())
	}

	/*tracer, closer := Init("client")
	defer closer.Close()*/

	helloTo := os.Args[1]
	for i := 0; i < 4; i++ {
		span := tracer.StartSpan("say-hello")
		span.SetTag("hello-to", helloTo)

		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		span.LogFields(
			log.String("event", "string-format"),
			log.String("value", fmt.Sprintf("%s%d", helloStr, i)),
		)

		println(helloStr)
		span.LogKV("event", "println")

		span.Finish()

	}
	time.Sleep(60 * time.Second)
	fmt.Printf("123")
	//closer.Close()
	fmt.Printf("456")
}

/*
func InitTrace(serviceName, reportAddr, sampleType string, rate float64) (tracer opentracing.Tracer,err error) {
	fmt.Println(serviceName, reportAddr, sampleType, rate)
	transport, err := zipkin.NewHTTPTransport(
		"http://127.0.0.1:9412/api/v1/spans",
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)

	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init zipkin: %v\n", err)
		return
	}

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)

	tracer, closer, err := cfg.New(serviceName,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init Jaeger: %v\n", err)
		return
	}
	_ = closer
	//opentracing.SetGlobalTracer(tracer)
	return
}
*/

func InitTrace(serviceName, reportAddr, sampleType string, rate float64) (tracer opentracing.Tracer, err error) {

	transport, err := zipkin.NewHTTPTransport(
		reportAddr,
		zipkin.HTTPBatchSize(16),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init zipkin: %v\n", err)
		return
	}

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: rate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)
	tracer, closer, err := cfg.New(serviceName,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init Jaeger: %v\n", err)
		return
	}

	_ = closer
	opentracing.SetGlobalTracer(tracer)
	return
}
