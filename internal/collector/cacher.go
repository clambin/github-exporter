package collector

import (
	"context"
	"fmt"
	"github.com/google/go-github/v58/github"
	"log/slog"
	"sync"
	"time"
)

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
		c.repoStats = uniqueRepoStats(stats)
		c.expiration = time.Now().Add(c.Lifetime)
	}

	return c.repoStats, err
}

func uniqueRepoStats(stats []RepoStats) []RepoStats {
	unique := make(map[string]RepoStats)
	for _, entry := range stats {
		unique[entry.GetFullName()] = entry
	}
	output := make([]RepoStats, 0, len(unique))
	for _, entry := range unique {
		output = append(output, entry)
	}
	return output
}

type RepoStats struct {
	*github.Repository
	PullRequestCount int
}

func (rs RepoStats) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", rs.GetFullName()),
		slog.Bool("archived", rs.GetArchived()),
		slog.Bool("fork", rs.GetFork()),
		slog.Bool("private", rs.GetPrivate()),
		slog.Int("stars", rs.GetStargazersCount()),
		slog.Int("issues", rs.GetOpenIssues()),
		slog.Int("pullRequests", rs.PullRequestCount),
	)
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

func (c *GitHubCache) queryAllRepoStats(ctx context.Context, ch chan repoStatResponse) {
	var wg sync.WaitGroup

	for _, user := range c.Users {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			c.queryUserRepoStats(ctx, ch, user)
		}(user)
	}

	for _, repo := range c.Repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			c.queryRepoStats(ctx, ch, repo)
		}(repo)
	}

	wg.Wait()
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
	var wg sync.WaitGroup

	for _, entry := range stats {
		wg.Add(1)
		go func(entry RepoStats) {
			defer wg.Done()
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
		}(entry)
	}

	wg.Wait()
	close(ch)
}
