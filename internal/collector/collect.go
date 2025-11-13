package collector

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/clambin/github-exporter/internal/stats/github"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	lastUpdate      time.Time
	Client          StatClient
	Logger          *slog.Logger
	Users           []string
	Repos           []string
	cache           []github.RepoStats
	Lifetime        time.Duration
	lock            sync.RWMutex
	IncludeArchived bool
}

type StatClient interface {
	GetRepoStats(context.Context, []string, []string) ([]github.RepoStats, error)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	defer func() { c.Logger.Debug("collected", "duration", time.Since(start)) }()

	repoStats, err := c.getStats()
	if err != nil {
		c.Logger.Error("failed to collect github statistics", "err", err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("github_exporter_error", "Error getting github statistics", nil, nil), err)
		return
	}

	for _, repoStat := range repoStats {
		c.Logger.Debug("repo found", "repo", repoStat.Name)

		if !c.IncludeArchived && repoStat.Archived {
			continue
		}

		archived := bool2string(repoStat.Archived)
		ch <- prometheus.MustNewConstMetric(metrics["stars"], prometheus.GaugeValue, float64(repoStat.Stars), repoStat.Name, archived)
		ch <- prometheus.MustNewConstMetric(metrics["forks"], prometheus.GaugeValue, float64(repoStat.Forks), repoStat.Name, archived)
		ch <- prometheus.MustNewConstMetric(metrics["issues"], prometheus.GaugeValue, float64(repoStat.Issues), repoStat.Name, archived)
		ch <- prometheus.MustNewConstMetric(metrics["pulls"], prometheus.GaugeValue, float64(repoStat.PullRequests), repoStat.Name, archived)
	}
}

func (c *Collector) getStats() ([]github.RepoStats, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if time.Since(c.lastUpdate) < c.Lifetime {
		return c.cache, nil
	}

	repoStats, err := c.Client.GetRepoStats(context.Background(), c.Users, c.Repos)
	if err == nil {
		c.cache = repoStats
		c.lastUpdate = time.Now()
	}

	return c.cache, err
}

func bool2string(val bool) string {
	booleans := map[bool]string{
		true:  "true",
		false: "false",
	}
	return booleans[val]
}
