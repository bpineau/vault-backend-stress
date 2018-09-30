package prometheus

import (
	"time"

	"github.com/bpineau/vault-backend-stress/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var (
	quantiles = map[float64]float64{
		0.5:  0.05,
		0.95: 0.005,
		0.99: 0.001,
	}

	promSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "success",
		Help: "Number of successfull requests",
	})

	promError = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "errors",
		Help: "Number of failed requests",
	})

	promTimes = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "timings",
		Help:       "Requests durations",
		Objectives: quantiles,
	})
)

func init() {
	prometheus.MustRegister(promSuccess, promError, promTimes)
}

// Sink collect metrics. Prom claims all exported functions are thread safe.
type Sink struct{}

// Observe receive an observation
func (n *Sink) Observe(status metrics.Status, duration time.Duration) {
	switch status {
	case metrics.Success:
		promSuccess.Inc()
	case metrics.Error:
		promError.Inc()
	}

	promTimes.Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

// Dump return metrics state
func (n *Sink) Dump() (*metrics.Point, error) {
	result := new(metrics.Point)
	result.Date = time.Now()

	metric := &dto.Metric{}
	err := promTimes.Write(metric)
	if err != nil {
		return nil, err
	}
	q := metric.GetSummary().GetQuantile()

	for _, v := range q {
		switch quant := v.GetQuantile(); quant {
		case 0.5:
			result.P50 = v.GetValue()
		case 0.95:
			result.P95 = v.GetValue()
		case 0.99:
			result.P99 = v.GetValue()
		}
	}

	err = promSuccess.Write(metric)
	if err != nil {
		return nil, err
	}
	result.SuccessCount = metric.GetCounter().GetValue()

	err = promError.Write(metric)
	if err != nil {
		return nil, err
	}
	result.ErrorsCount = metric.GetCounter().GetValue()

	return result, nil
}
