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

type Cacher struct {
	Client          GitHubClient
	Users           []string
	Repos           []string
	IncludeArchived bool
	Lifetime        time.Duration
	content         []RepoStats
	created         time.Time
	lock            sync.Mutex
}

func (c *Cacher) Get() ([]RepoStats, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.created.IsZero() && time.Now().Before(c.created.Add(c.Lifetime)) {
		return c.content, nil
	}

	ch := make(chan RepoStatResponse)
	go c.GetStats(ch)

	c.content = make([]RepoStats, 0)
	var err error
	for resp := range ch {
		if resp.Err != nil {
			err = resp.Err
			continue
		}
		c.content = append(c.content, resp.Stats)
	}
	c.created = time.Now()

	return c.content, err
}

type RepoStatResponse struct {
	Stats RepoStats
	Err   error
}

type RepoStats struct {
	*github.Repository
	PullRequestCount int
}

func (c *Cacher) GetStats(ch chan RepoStatResponse) {
	ctx := context.Background()
	repos := make(chan RepoStatResponse)
	go c.queryAllRepoStats(ctx, repos)

	var stats []RepoStats
	var err error

	for resp := range repos {
		if resp.Err != nil {
			err = resp.Err
		}
		if resp.Stats.Repository.GetArchived() && !c.IncludeArchived {
			continue
		}
		stats = append(stats, resp.Stats)
	}
	if err != nil {
		ch <- RepoStatResponse{Err: fmt.Errorf("get repo stats: %w", err)}
		close(ch)
		return
	}

	c.addPullRequests(ctx, stats, ch)
}

const maxParallel = 100

func (c *Cacher) queryAllRepoStats(ctx context.Context, ch chan RepoStatResponse) {
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

func (c *Cacher) queryUserRepoStats(ctx context.Context, ch chan RepoStatResponse, user string) {
	userRepos, err := c.Client.GetUserRepos(ctx, user)
	if err != nil {
		ch <- RepoStatResponse{Err: err}
		return
	}
	for _, userRepo := range userRepos {
		ch <- RepoStatResponse{Stats: RepoStats{Repository: userRepo}}
	}
}

func (c *Cacher) queryRepoStats(ctx context.Context, ch chan RepoStatResponse, repoName string) {
	repo, err := c.Client.GetRepo(ctx, repoName)
	ch <- RepoStatResponse{
		Stats: RepoStats{Repository: repo},
		Err:   err,
	}
}

func (c *Cacher) addPullRequests(ctx context.Context, stats []RepoStats, ch chan RepoStatResponse) {
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

			ch <- RepoStatResponse{
				Stats: RepoStats{Repository: entry.Repository, PullRequestCount: pullRequestCount},
				Err:   err,
			}
			parallel.Release(1)
		}(entry)
	}
	_ = parallel.Acquire(ctx, maxParallel)
	close(ch)
}
