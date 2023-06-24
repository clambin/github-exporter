package collector

import (
	"context"
	github2 "github.com/clambin/github-exporter/pkg/github"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

var _ prometheus.Collector = &Collector{}

//go:generate mockery --name GitHubClient
type GitHubClient interface {
	GetUserRepos(context.Context, string) ([]github2.Repo, error)
	GetRepo(context.Context, string) (github2.Repo, error)
	GetPullRequests(context.Context, string) ([]github2.PullRequest, error)
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
	ch2 := make(chan repoStatResponse)
	go c.getStats(ch2)

	for entry := range ch2 {
		if entry.err != nil {
			slog.Error("failed to collect github statistics", "err", entry.err)
			ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("github_monitor_error", "Error getting github statistics", nil, nil), entry.err)
			//return
			continue
		}
		archived := bool2string(entry.stats.Repo.Archived)
		fork := bool2string(entry.stats.Repo.Fork)
		private := bool2string(entry.stats.Repo.Private)

		ch <- prometheus.MustNewConstMetric(metrics["stars"], prometheus.GaugeValue, float64(entry.stats.Repo.StargazersCount), entry.stats.Repo.FullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["forks"], prometheus.GaugeValue, float64(entry.stats.Repo.ForksCount), entry.stats.Repo.FullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["issues"], prometheus.GaugeValue, float64(entry.stats.Repo.OpenIssuesCount), entry.stats.Repo.FullName, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["pulls"], prometheus.GaugeValue, float64(entry.stats.pullRequestCount), entry.stats.Repo.FullName, archived, fork, private)
	}
}

func bool2string(val bool) string {
	booleans := map[bool]string{
		true:  "true",
		false: "false",
	}
	return booleans[val]
}
