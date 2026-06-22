package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "proxy_active_connections",
			Help: "Current number of active proxy connections",
		},
	)

	BytesInTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_bytes_in_total",
			Help: "Total bytes received from clients",
		},
		[]string{"downstream_az"},
	)

	BytesOutTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_bytes_out_total",
			Help: "Total bytes sent to clients (from downstream)",
		},
		[]string{"downstream_az"},
	)
)
