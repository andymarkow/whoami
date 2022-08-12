package telemetry

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func Init(version string) {
	promBuildInfo := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "whoami",
			Subsystem: "build",
			Name:      "info",
			Help:      "",
			ConstLabels: map[string]string{
				"version": version,
			},
		})

	promRuntimeInfo := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "whoami",
			Subsystem: "runtime",
			Name:      "info",
			Help:      "",
			ConstLabels: map[string]string{
				"go_version": runtime.Version(),
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
			},
		})

	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	promBuildInfo.Set(1)
	promRuntimeInfo.Set(1)
}
