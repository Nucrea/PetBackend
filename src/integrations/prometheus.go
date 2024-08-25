package integrations

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	reg            *prometheus.Registry
	rpsGauge       prometheus.Gauge
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
	rpsGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
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
	reg.MustRegister(rpsGauge, avgReqTimeHist, panicsHist)

	return &Prometheus{
		panicsHist:     panicsHist,
		avgReqTimeHist: avgReqTimeHist,
		rpsGauge:       rpsGauge,
		reg:            reg,
	}
}

func (p *Prometheus) GetRequestHandler() http.Handler {
	return promhttp.HandlerFor(p.reg, promhttp.HandlerOpts{Registry: p.reg})
}

func (p *Prometheus) RequestInc() {
	p.rpsGauge.Inc()
}

func (p *Prometheus) RequestDec() {
	p.rpsGauge.Dec()
}

func (p *Prometheus) AddRequestTime(reqTime float64) {
	p.avgReqTimeHist.Observe(reqTime)
}

func (p *Prometheus) AddPanic() {
	p.panicsHist.Observe(1)
}
