package collector

import (
	"context"
	"fmt"
	"github.com/clambin/github-exporter/pkg/github"
	"golang.org/x/sync/semaphore"
)

type repoStats struct {
	github.Repo
	pullRequestCount int
}

type repoStatResponse struct {
	stats repoStats
	err   error
}

const maxParallel = 20

func (c Collector) getStats() ([]repoStats, error) {
	ctx := context.Background()
	ch := make(chan repoStatResponse)
	go c.queryAllRepoStats(ctx, ch)

	var err error
	var stats []repoStats
	for resp := range ch {
		if resp.err != nil {
			err = fmt.Errorf("get repo stats: %w", resp.err)
		}
		if resp.stats.Repo.Archived && !c.IncludeArchived {
			continue
		}
		stats = append(stats, resp.stats)
	}
	if err != nil {
		return nil, err
	}

	ch = make(chan repoStatResponse)
	go c.getPRs(ctx, stats, ch)
	stats = make([]repoStats, 0, len(stats))
	for resp := range ch {
		if resp.err != nil {
			err = fmt.Errorf("get repo pr's: %w", resp.err)
		}
		stats = append(stats, resp.stats)
	}
	return stats, err
}

func (c Collector) queryAllRepoStats(ctx context.Context, ch chan repoStatResponse) {
	parallel := semaphore.NewWeighted(maxParallel)

	for _, user := range c.Users {
		_ = parallel.Acquire(ctx, 1)
		go func(user string) {
			c.queryUserRepoStats(ctx, ch, user)
			parallel.Release(1)
		}(user)
	}

	for _, repo := range c.Repos {
		_ = parallel.Acquire(ctx, 1)
		go func(repo string) {
			c.queryRepoStats(ctx, ch, repo)
			parallel.Release(1)
		}(repo)
	}

	_ = parallel.Acquire(ctx, maxParallel)
	close(ch)
}

func (c Collector) queryUserRepoStats(ctx context.Context, ch chan repoStatResponse, user string) {
	userRepos, err := c.Client.GetUserRepos(ctx, user)
	if err != nil {
		ch <- repoStatResponse{err: err}
		return
	}
	for _, userRepo := range userRepos {
		ch <- repoStatResponse{stats: repoStats{Repo: userRepo}}
	}
}

func (c Collector) queryRepoStats(ctx context.Context, ch chan repoStatResponse, repoName string) {
	repo, err := c.Client.GetRepo(ctx, repoName)
	ch <- repoStatResponse{
		stats: repoStats{Repo: repo},
		err:   err,
	}
}

func (c Collector) getPRs(ctx context.Context, stats []repoStats, ch chan repoStatResponse) {
	parallel := semaphore.NewWeighted(maxParallel)

	for _, entry := range stats {
		_ = parallel.Acquire(ctx, 1)
		go func(entry repoStats) {
			pullRequests, err := c.Client.GetPullRequests(ctx, entry.Repo.FullName)
			if err != nil {
				ch <- repoStatResponse{err: err}
				return
			}
			entry.Repo.OpenIssuesCount -= len(pullRequests)
			ch <- repoStatResponse{stats: repoStats{
				Repo:             entry.Repo,
				pullRequestCount: len(pullRequests),
			}}
			parallel.Release(1)
		}(entry)
	}
	_ = parallel.Acquire(ctx, maxParallel)
	close(ch)
}
