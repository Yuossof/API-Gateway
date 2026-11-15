package matrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_requests_total",
		Help: "Total number of requests received by the gateway",
	},
	[]string{"method"},
)

func Handler() http.Handler {
	return promhttp.Handler()
}

func StartMatricsServer() {
	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/matrics", Handler())

		log.Println("matrics server is running on :2112")

		if err := http.ListenAndServe(":2112", metricsMux); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()
}
