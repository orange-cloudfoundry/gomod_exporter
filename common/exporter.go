package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

// Metrics - hold metrics and initialized registry
type Metrics struct {
	Info       *prometheus.GaugeVec
	Deprecated *prometheus.GaugeVec
	Replaced   *prometheus.GaugeVec
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
		Replaced: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: ns,
				Name:      "replaced",
				Help:      "Give information about module replacements",
			},
			[]string{"module", "dependency", "type", "replacement", "version"},
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

	if err := res.Registry.Register(res.Info); err != nil {
		log.Errorf("unable to register info metric: %s", err)
	}
	if err := res.Registry.Register(res.Deprecated); err != nil {
		log.Errorf("unable to register deprecated metric: %s", err)
	}
	if err := res.Registry.Register(res.Replaced); err != nil {
		log.Errorf("unable to register replaced metric: %s", err)
	}
	if err := res.Registry.Register(res.Status); err != nil {
		log.Errorf("unable to register status metric: %s", err)
	}
	if err := res.Registry.Register(res.Duration); err != nil {
		log.Errorf("unable to register duration metric: %s", err)
	}
	return res
}
