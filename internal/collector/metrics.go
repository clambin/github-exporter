package collector

import "github.com/prometheus/client_golang/prometheus"

var metrics = map[string]*prometheus.Desc{
	"stars": prometheus.NewDesc(
		prometheus.BuildFQName("github", "exporter", "stars"),
		"Total number of stars",
		[]string{"repo", "archived"},
		nil,
	),
	"issues": prometheus.NewDesc(
		prometheus.BuildFQName("github", "exporter", "issues"),
		"Total number of open issues",
		[]string{"repo", "archived"},
		nil,
	),
	"pulls": prometheus.NewDesc(
		prometheus.BuildFQName("github", "exporter", "pulls"),
		"Total number of open pull requests",
		[]string{"repo", "archived"},
		nil,
	),
	"forks": prometheus.NewDesc(
		prometheus.BuildFQName("github", "exporter", "forks"),
		"Total number of forks",
		[]string{"repo", "archived"},
		nil,
	),
}
