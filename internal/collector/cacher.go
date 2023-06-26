package collector

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

//go:generate mockery --name GitHubClient
type GitHubClient interface {
	GetUserRepos(context.Context, string) ([]*github.Repository, error)
	GetRepo(context.Context, string) (*github.Repository, error)
	GetPullRequests(context.Context, string) ([]*github.PullRequest, error)
}

type GitHubCache struct {
	Client          GitHubClient
	Users           []string
	Repos           []string
	IncludeArchived bool
	Lifetime        time.Duration
	repoStats       []RepoStats
	expiration      time.Time
	lock            sync.Mutex
}

func (c *GitHubCache) Get() ([]RepoStats, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if time.Now().Before(c.expiration) {
		return c.repoStats, nil
	}

	stats, err := c.getStats()
	if err == nil {
		c.repoStats = stats
		c.expiration = time.Now().Add(c.Lifetime)
	}

	return c.repoStats, err
}

type RepoStats struct {
	*github.Repository
	PullRequestCount int
}

type repoStatResponse struct {
	Stats RepoStats
	Err   error
}

func (c *GitHubCache) getStats() ([]RepoStats, error) {
	ctx := context.Background()
	stats, err := c.getAllRepoStats(ctx)
	if err != nil {
		return nil, err
	}
	return c.addPullRequests(ctx, stats)
}

func (c *GitHubCache) getAllRepoStats(ctx context.Context) ([]RepoStats, error) {
	ch := make(chan repoStatResponse)
	go c.queryAllRepoStats(ctx, ch)

	var stats []RepoStats
	var err error

	for resp := range ch {
		if resp.Err != nil {
			err = resp.Err
		}
		if resp.Stats.Repository.GetArchived() && !c.IncludeArchived {
			continue
		}
		stats = append(stats, resp.Stats)
	}
	return stats, err
}

const maxParallel = 25

func (c *GitHubCache) queryAllRepoStats(ctx context.Context, ch chan repoStatResponse) {
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

func (c *GitHubCache) queryUserRepoStats(ctx context.Context, ch chan repoStatResponse, user string) {
	userRepos, err := c.Client.GetUserRepos(ctx, user)
	if err != nil {
		ch <- repoStatResponse{Err: err}
		return
	}
	for _, userRepo := range userRepos {
		ch <- repoStatResponse{Stats: RepoStats{Repository: userRepo}}
	}
}

func (c *GitHubCache) queryRepoStats(ctx context.Context, ch chan repoStatResponse, repoName string) {
	repo, err := c.Client.GetRepo(ctx, repoName)
	ch <- repoStatResponse{
		Stats: RepoStats{Repository: repo},
		Err:   err,
	}
}

func (c *GitHubCache) addPullRequests(ctx context.Context, stats []RepoStats) ([]RepoStats, error) {
	ch := make(chan repoStatResponse)
	go c.getPullRequests(ctx, stats, ch)

	var err error
	result := make([]RepoStats, 0, len(stats))
	for stat := range ch {
		if stat.Err != nil {
			err = stat.Err
		}
		result = append(result, stat.Stats)
	}
	return result, err
}

func (c *GitHubCache) getPullRequests(ctx context.Context, stats []RepoStats, ch chan repoStatResponse) {
	parallel := semaphore.NewWeighted(maxParallel)

	for _, entry := range stats {
		_ = parallel.Acquire(ctx, 1)
		go func(entry RepoStats) {
			fullName := entry.Repository.GetFullName()

			var pullRequestCount int
			pullRequests, err := c.Client.GetPullRequests(ctx, fullName)

			if err == nil {
				pullRequestCount = len(pullRequests)
				*entry.Repository.OpenIssuesCount -= pullRequestCount
			} else {
				err = fmt.Errorf("get pr stats: %w", err)
			}

			ch <- repoStatResponse{
				Stats: RepoStats{Repository: entry.Repository, PullRequestCount: pullRequestCount},
				Err:   err,
			}
			parallel.Release(1)
		}(entry)
	}
	_ = parallel.Acquire(ctx, maxParallel)
	close(ch)
}
