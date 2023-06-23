package collector

import (
	"context"
	"github.com/clambin/github-exporter/internal/github"
	"golang.org/x/sync/semaphore"
)

type repoStats struct {
	repo github.Repo
	prs  int
}

func (c Collector) getStats() ([]repoStats, error) {
	ch := make(chan repoStatResponse)
	go c.queryAllRepoStats(ch)

	var stats []repoStats
	for resp := range ch {
		if resp.err != nil {
			return nil, resp.err
		}
		if resp.stats.repo.Archived && !c.IncludeArchived {
			continue
		}
		stats = append(stats, resp.stats)
	}

	ch = make(chan repoStatResponse)
	go c.getPRs(stats, ch)
	stats = make([]repoStats, 0, len(stats))
	for resp := range ch {
		if resp.err != nil {
			return nil, resp.err
		}
		stats = append(stats, resp.stats)
	}
	return stats, nil
}

type repoStatResponse struct {
	stats repoStats
	err   error
}

const maxParallel = 20

func (c Collector) queryAllRepoStats(ch chan repoStatResponse) {
	parallel := semaphore.NewWeighted(maxParallel)
	ctx := context.Background()

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
		ch <- repoStatResponse{stats: repoStats{repo: userRepo}}
	}
}

func (c Collector) queryRepoStats(ctx context.Context, ch chan repoStatResponse, repoName string) {
	repo, err := c.Client.GetRepo(ctx, repoName)
	ch <- repoStatResponse{
		stats: repoStats{repo: repo},
		err:   err,
	}
}

func (c Collector) getPRs(stats []repoStats, ch chan repoStatResponse) {
	parallel := semaphore.NewWeighted(maxParallel)
	ctx := context.Background()

	for _, entry := range stats {
		_ = parallel.Acquire(ctx, 1)
		go func(entry repoStats) {
			prs, err := c.Client.GetPullRequests(ctx, entry.repo.FullName)
			if err != nil {
				ch <- repoStatResponse{err: err}
				return
			}
			entry.repo.OpenIssuesCount -= len(prs)
			ch <- repoStatResponse{stats: repoStats{
				repo: entry.repo,
				prs:  len(prs),
			}}
			parallel.Release(1)
		}(entry)
	}
	_ = parallel.Acquire(ctx, maxParallel)
	close(ch)
}
