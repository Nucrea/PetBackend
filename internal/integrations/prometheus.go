package integrations

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Counter interface {
	Inc()
}

type Gauge interface {
	Set(float64)
	Inc()
	Dec()
}

type Histogram interface {
	Observe(float64)
}

type Metrics struct {
	registry   *prometheus.Registry
	registerer prometheus.Registerer
}

func NewMetrics(prefix string) *Metrics {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWithPrefix(prefix, registry)

	registerer.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &Metrics{
		registry:   registry,
		registerer: registerer,
	}
}

func (m *Metrics) NewCounter(name, description string) Counter {
	collector := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *Metrics) NewGauge(name, description string) Gauge {
	collector := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *Metrics) NewHistogram(name, description string) Histogram {
	collector := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *Metrics) HttpHandler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registerer})
}
