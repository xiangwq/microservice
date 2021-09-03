package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

type Metrics struct {
	requestCounter *prometheus.CounterVec
	codeCounter    *prometheus.CounterVec
	latencySummary *prometheus.SummaryVec
}

func NewServerMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "microservice_server_request_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method"}),
		codeCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "microservice_server_handled_code_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "grpc_code"}),
		latencySummary: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "microservice_proc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

func (m *Metrics) IncrRequest(ctx context.Context, serviceName, methodName string) {
	m.requestCounter.WithLabelValues(serviceName, methodName).Inc()
}

func (m *Metrics) IncrCode(ctx context.Context, serviceName, methodName string, err error) {
	st, _ := status.FromError(err)
	m.codeCounter.WithLabelValues(serviceName, methodName, st.Code().String()).Inc()
}

func (m *Metrics) Latency(ctx context.Context, serviceName, methodName string, us int64) {

	m.latencySummary.WithLabelValues(serviceName, methodName).Observe(float64(us))
}

//生成server metrics实例
func NewRpcMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "microservice_rpc_call_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method"}),
		codeCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "microservice_rpc_code_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "grpc_code"}),
		latencySummary: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "microservice_rpc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}