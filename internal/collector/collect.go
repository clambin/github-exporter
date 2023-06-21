package collector

import (
	"context"
	"fmt"
	"github.com/clambin/github-exporter/internal/github"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	Client          github.Client
	Users           []string
	Repos           []string
	IncludeArchived bool
}

var metrics = map[string]*prometheus.Desc{
	"stars": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "stars"),
		"Total number of stars",
		[]string{"repo", "user", "archived", "fork", "private"},
		nil,
	),
	"issues": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "issues"),
		"Total number of open issues",
		[]string{"repo", "user", "archived", "fork", "private"},
		nil,
	),
	"pulls": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "pulls"),
		"Total number of open pull requests",
		[]string{"repo", "user", "archived", "fork", "private"},
		nil,
	),
	"forks": prometheus.NewDesc(
		prometheus.BuildFQName("github", "monitor", "forks"),
		"Total number of forks",
		[]string{"repo", "user", "archived", "fork", "private"},
		nil,
	),
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	repos, err := c.getAllRepos(ctx)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("github_monitor_error", "Error getting github statistics", nil, nil), err)
		return
	}

	for _, repo := range repos {
		if !c.IncludeArchived && repo.Archived {
			continue
		}
		user, archived, fork, private := analyzeRepo(repo)

		ch <- prometheus.MustNewConstMetric(metrics["stars"], prometheus.GaugeValue, float64(repo.StargazersCount), repo.FullName, user, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["forks"], prometheus.GaugeValue, float64(repo.ForksCount), repo.FullName, user, archived, fork, private)

		pullRequests, err := c.Client.GetPullRequests(ctx, repo.FullName)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("github_monitor_error", "Error getting github statistics", nil, nil), err)
			return
		}

		ch <- prometheus.MustNewConstMetric(metrics["pulls"], prometheus.GaugeValue, float64(len(pullRequests)), repo.FullName, user, archived, fork, private)
		ch <- prometheus.MustNewConstMetric(metrics["issues"], prometheus.GaugeValue, float64(repo.OpenIssuesCount-len(pullRequests)), repo.FullName, user, archived, fork, private)
	}
}

func (c Collector) getAllRepos(ctx context.Context) ([]github.Repo, error) {
	var repos []github.Repo

	for _, user := range c.Users {
		userRepos, err := c.Client.GetUserRepos(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to get repos for '%s': %w", user, err)
		}
		repos = append(repos, userRepos...)
	}
	for _, repo := range c.Repos {
		userRepo, err := c.Client.GetRepo(ctx, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get repo '%s': %w", repo, err)
		}
		repos = append(repos, userRepo)

	}

	return repos, nil
}

func analyzeRepo(repo github.Repo) (user, archived, fork, private string) {
	if parts := strings.Split(repo.FullName, "/"); len(parts) == 2 {
		user = parts[0]
	} else {
		user = "unknown"
	}
	return user, bool2string(repo.Archived), bool2string(repo.Fork), bool2string(repo.Private)
}

func bool2string(val bool) string {
	if val {
		return "true"
	}
	return "false"
}
