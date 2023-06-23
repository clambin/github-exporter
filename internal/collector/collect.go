package collector

import (
	"context"
	"github.com/clambin/github-exporter/internal/github"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

var _ prometheus.Collector = &Collector{}

//go:generate mockery --name GitHubClient
type GitHubClient interface {
	GetUserRepos(context.Context, string) ([]github.Repo, error)
	GetRepo(context.Context, string) (github.Repo, error)
	GetPullRequests(context.Context, string) ([]github.PullRequest, error)
}

type Collector struct {
	Client          GitHubClient
	Users           []string
	Repos           []string
	IncludeArchived bool
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

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.getStats()
	if err != nil {
		slog.Error("failed to collect github statistics", "err", err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("github_monitor_error", "Error getting github statistics", nil, nil), err)
		return
	}

	for _, entry := range stats {
		archived := bool2string(entry.repo.Archived)
		fork := bool2string(entry.repo.Fork)
		private := bool2string(entry.repo.Private)

		ch <- prometheus.MustNewConstMetric(metrics["stars"], prometheus.GaugeValue, float64(entry.repo.StargazersCount), entry.repo.FullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["forks"], prometheus.GaugeValue, float64(entry.repo.ForksCount), entry.repo.FullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["issues"], prometheus.GaugeValue, float64(entry.repo.OpenIssuesCount), entry.repo.FullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["pulls"], prometheus.GaugeValue, float64(entry.prs), entry.repo.FullName, archived, fork, private)
	}
}

func bool2string(val bool) string {
	booleans := map[bool]string{
		true:  "true",
		false: "false",
	}
	return booleans[val]
}
