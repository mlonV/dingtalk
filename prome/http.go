package prome

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var PromeRegister *prometheus.Registry

func PromeHTTPFunc() gin.HandlerFunc {
	PromeRegister = prometheus.NewRegistry()

	h := promhttp.HandlerFor(PromeRegister, promhttp.HandlerOpts{Registry: PromeRegister})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// // 进程监控使用 For example
// func ProcessGauge() prometheus.Gauge {
// 	return prometheus.NewGauge(prometheus.GaugeOpts{
// 		Namespace: "our_company",
// 		Subsystem: "blob_storage",
// 		Name:      "ops_queued",
// 		Help:      "Number of blob storage operations waiting to be processed.",
// 	})

// }
