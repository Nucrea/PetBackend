package integrations

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	reg            *prometheus.Registry
	rpsCounter     prometheus.Counter
	avgReqTimeHist prometheus.Histogram
	panicsHist     prometheus.Histogram
}

func NewPrometheus() *Prometheus {
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// errorsCounter := prometheus.NewCounter(
	// 	prometheus.CounterOpts{
	// 		Name: "backend_errors_count",
	// 		Help: "Summary errors count",
	// 	},
	// )
	rpsCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "backend_requests_per_second",
			Help: "Requests per second metric",
		},
	)
	avgReqTimeHist := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "backend_requests_average_time",
			Help: "Average time of requests",
		},
	)
	panicsHist := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "backend_panics",
			Help: "Panics histogram metric",
		},
	)
	reg.MustRegister(rpsCounter, avgReqTimeHist, panicsHist)

	return &Prometheus{
		panicsHist:     panicsHist,
		avgReqTimeHist: avgReqTimeHist,
		rpsCounter:     rpsCounter,
		reg:            reg,
	}
}

func (p *Prometheus) GetRequestHandler() http.Handler {
	return promhttp.HandlerFor(p.reg, promhttp.HandlerOpts{Registry: p.reg})
}

func (p *Prometheus) RequestInc() {
	p.rpsCounter.Inc()
}

func (p *Prometheus) RequestDec() {
	// p.rpsGauge.Dec()
}

func (p *Prometheus) AddRequestTime(reqTime float64) {
	p.avgReqTimeHist.Observe(reqTime)
}

func (p *Prometheus) AddPanic() {
	p.panicsHist.Observe(1)
}
