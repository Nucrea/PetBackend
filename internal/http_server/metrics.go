package httpserver

import (
	"backend/internal/integrations"
)

func NewServerMetrics(p *integrations.Metrics) *ServerMetrics {
	errors5xxCounter := p.NewCounter("server_responses_5xx", "5xx responses counter")
	errors4xxCounter := p.NewCounter("server_responses_4xx", "4xx responses count")
	requestsCounter := p.NewCounter("server_requests_total", "requests counter")
	avgReqTimeHist := p.NewHistogram("server_requests_time", "requests time histogram")
	panicsHist := p.NewHistogram("server_panics", "panics histogram metric")

	return &ServerMetrics{
		rpsCounter:       requestsCounter,
		avgReqTimeHist:   avgReqTimeHist,
		panicsHist:       panicsHist,
		errors4xxCounter: errors4xxCounter,
		errors5xxCounter: errors5xxCounter,
	}
}

type ServerMetrics struct {
	rpsCounter       integrations.Counter
	avgReqTimeHist   integrations.Histogram
	panicsHist       integrations.Histogram
	errors4xxCounter integrations.Counter
	errors5xxCounter integrations.Counter
}

func (b *ServerMetrics) AddRequest() {
	b.rpsCounter.Inc()
}

func (b *ServerMetrics) AddRequestTime(reqTime float64) {
	b.avgReqTimeHist.Observe(reqTime)
}

func (b *ServerMetrics) AddPanic() {
	b.panicsHist.Observe(1)
}

func (b *ServerMetrics) Add4xxError() {
	b.errors4xxCounter.Inc()
}

func (b *ServerMetrics) Add5xxError() {
	b.errors5xxCounter.Inc()
}
