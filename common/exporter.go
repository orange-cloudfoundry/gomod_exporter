package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics - hold metrics and initialized registry
type Metrics struct {
	Info       *prometheus.GaugeVec
	Deprecated *prometheus.GaugeVec
	Status     *prometheus.GaugeVec
	Duration   prometheus.Gauge
	Registry   *prometheus.Registry
}

// NewMetrics - create Metrics object
func NewMetrics(ns string) *Metrics {
	res := &Metrics{
		Info: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "info",
				Help:      "Informations about given repository, value always 1",
			},
			[]string{"module", "goversion"},
		),
		Deprecated: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "deprecated",
				Help:      "Number of days since given dependency of repository is out-of-date",
			},
			[]string{"module", "dependency", "type", "current", "latest"},
		),
		Status: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "status",
				Help:      "Status of last analysis of given repository, 0 for error",
			},
			[]string{"repository"},
		),
		Duration: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "duration",
				Help:      "Duration of last analysis in second",
			},
		),
		Registry: prometheus.NewRegistry(),
	}

	res.Registry.Register(res.Info)
	res.Registry.Register(res.Deprecated)
	res.Registry.Register(res.Status)
	res.Registry.Register(res.Duration)
	return res
}
