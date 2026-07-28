package stats

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"strings"
	"sync"
	"time"

	"codeberg.org/clambin/go-common/set"
	"github.com/clambin/github-exporter/internal/stats/github"
)

type Client struct {
	GitHubClient
	Logger *slog.Logger
}

type GitHubClient interface {
	GetUserRepoNames(context.Context, string) ([]string, error)
	GetRepoStats(context.Context, string, string) (github.RepoStats, error)
	GetPullRequestCount(context.Context, string, string) (int, error)
}

func (c Client) GetRepoStats(ctx context.Context, users []string, repos []string) ([]github.RepoStats, error) {
	var wg sync.WaitGroup
	type result struct {
		stats github.RepoStats
		err   error
	}
	ch := make(chan result)

	var count int
	for repoName, err := range c.uniqueRepoNames(ctx, users, repos) {
		if err != nil {
			return nil, err
		}
		c.Logger.Debug("repo found", "repo", repoName)
		count++
		wg.Go(func() {
			stats, err := c.getStats(ctx, repoName)
			ch <- result{stats: stats, err: err}
		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	errs := make([]error, 0, count)
	stats := make([]github.RepoStats, 0, count)
	for r := range ch {
		if r.err == nil {
			stats = append(stats, r.stats)
		}
		errs = append(errs, r.err)
	}
	return stats, errors.Join(errs...)
}

func (c Client) uniqueRepoNames(ctx context.Context, users []string, repos []string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		uniqueRepoNames := set.New(repos...)
		for _, user := range users {
			userRepos, err := c.GetUserRepoNames(ctx, user)
			if err != nil {
				yield("", fmt.Errorf("get repos for user %s: %w", user, err))
				return
			}
			for _, userRepo := range userRepos {
				if !uniqueRepoNames.Contains(userRepo) {
					if !yield(userRepo, nil) {
						return
					}
					uniqueRepoNames.Add(userRepo)
				}
			}
		}
	}
}

func (c Client) getStats(ctx context.Context, repo string) (github.RepoStats, error) {
	start := time.Now()
	defer func() {
		c.Logger.Debug("got repo stats", "repo", repo, "duration", time.Since(start))
	}()

	user, repo, err := splitFullName(repo)
	if err != nil {
		return github.RepoStats{}, err
	}

	repoStats, err := c.GitHubClient.GetRepoStats(ctx, user, repo)
	if err != nil {
		return repoStats, err
	}
	repoStats.PullRequests, err = c.GetPullRequestCount(ctx, user, repo)
	if err != nil {
		return repoStats, err
	}
	repoStats.Issues -= repoStats.PullRequests
	return repoStats, nil
}

func splitFullName(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo name: %s", repo)
	}
	return parts[0], parts[1], nil
}
