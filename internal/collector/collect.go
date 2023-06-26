package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	GitHubCache
}

var metrics = map[string]*prometheus.Desc{
	"stars": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "stars"),
		"Total number of stars",
		[]string{"repo", "archived", "fork", "private"},
		nil,
	),
	"issues": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "issues"),
		"Total number of open issues",
		[]string{"repo", "archived", "fork", "private"},
		nil,
	),
	"pulls": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "pulls"),
		"Total number of open pull requests",
		[]string{"repo", "archived", "fork", "private"},
		nil,
	),
	"forks": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "forks"),
		"Total number of forks",
		[]string{"repo", "archived", "fork", "private"},
		nil,
	),
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.GitHubCache.Get()

	if err != nil {
		slog.Error("failed to collect github statistics", "err", err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("github_monitor_error", "Error getting github statistics", nil, nil), err)
		return
	}

	for _, entry := range stats {
		fullName := entry.Repository.GetFullName()
		archived := bool2string(entry.Repository.GetArchived())
		fork := bool2string(entry.Repository.GetFork())
		private := bool2string(entry.Repository.GetPrivate())

		ch <- prometheus.MustNewConstMetric(metrics["stars"], prometheus.GaugeValue, float64(entry.Repository.GetStargazersCount()), fullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["forks"], prometheus.GaugeValue, float64(entry.Repository.GetForksCount()), fullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["issues"], prometheus.GaugeValue, float64(entry.Repository.GetOpenIssuesCount()), fullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["pulls"], prometheus.GaugeValue, float64(entry.PullRequestCount), fullName, archived, fork, private)
	}
}

func bool2string(val bool) string {
	booleans := map[bool]string{
		true:  "true",
		false: "false",
	}
	return booleans[val]
}
