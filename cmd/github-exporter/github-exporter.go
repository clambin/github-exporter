package main

import (
	"github.com/clambin/github-exporter/internal/collector"
	"github.com/clambin/github-exporter/internal/github"
	"github.com/clambin/go-common/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"time"
)

func main() {
	user := os.Getenv("GITHUB_EXPORTER_USER")
	token := os.Getenv("GITHUB_EXPORTER_TOKEN")

	c := collector.Collector{
		Users: []string{user},
		Client: github.Client{
			HTTPClient: &http.Client{
				Transport: httpclient.NewRoundTripper(httpclient.WithCache(httpclient.DefaultCacheTable, time.Hour, 24*time.Hour)),
				Timeout:   10 * time.Second,
			},
			Token: token,
		},
	}

	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":9090", nil)
}
