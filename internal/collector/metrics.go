package collector

import "github.com/prometheus/client_golang/prometheus"

var metrics = map[string]*prometheus.Desc{
	"stars": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "stars"),
		"Total number of stars",
		[]string{"repo", "archived"},
		nil,
	),
	"issues": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "issues"),
		"Total number of open issues",
		[]string{"repo", "archived"},
		nil,
	),
	"pulls": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "pulls"),
		"Total number of open pull requests",
		[]string{"repo", "archived"},
		nil,
	),
	"forks": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "forks"),
		"Total number of forks",
		[]string{"repo", "archived"},
		nil,
	),
}
