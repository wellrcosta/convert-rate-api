package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Requests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currency_converter_requests_total",
			Help: "Total number of requests received",
		},
	)

	Conversions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "currency_converter_conversions_total",
			Help: "Total number of currency conversions",
		},
		[]string{"from", "to"},
	)

	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currency_converter_cache_hits_total",
			Help: "Number of times conversion rate was served from cache",
		},
	)

	CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currency_converter_cache_misses_total",
			Help: "Number of times conversion rate was fetched from API",
		},
	)

	Errors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "currency_converter_errors_total",
			Help: "Number of errors encountered during processing",
		},
	)
)

func init() {
	prometheus.MustRegister(Requests, Conversions, CacheHits, CacheMisses, Errors)
}

// PrometheusHandler returns the HTTP handler for Prometheus metrics.
func PrometheusHandler() func(ctx *gin.Context) {
	h := promhttp.Handler()
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
