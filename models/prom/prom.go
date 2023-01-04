package prom

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	v1HttpRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "v1_http_requests",
		Help: "The total number of http requests made to the v1 api",
	}, []string{"code", "path"})
)

func Init() http.Handler {
	prometheus.MustRegister(v1HttpRequests)
	return promhttp.Handler()
}

func OnNewV1Request(path string, statusCode int) {
	v1HttpRequests.WithLabelValues(strconv.Itoa(statusCode), path).Inc()
}
