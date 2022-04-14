package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "civoapp_online_users",
	Help: "Online users",
	ConstLabels: map[string]string{
		"app": "civoapp",
	},
})

var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "civoapp_http_requests_total",
	Help: "Count of all HTTP requests for civoapp",
}, []string{})

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "civoapp_http_request_duration",
	Help: "Duration in seconds of all HTTP requests",
}, []string{"handler"})

func main() {
	r := prometheus.NewRegistry()
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpDuration)

	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(5000)))
		}
	}()

	homePageFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(9)) * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from inside a Civo Kubernetes Cluster"))
	})

	errorPageFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(7)) * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("You shouldn't be here"))
	})

	homePage := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "homePage"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, homePageFunc),
	)

	errorPage := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "errorPage"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, errorPageFunc),
	)

	http.Handle("/", homePage)
	http.Handle("/error", errorPage)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}