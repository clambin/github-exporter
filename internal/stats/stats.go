package stats

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"strings"
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
	var p parallel[github.RepoStats]
	for repoName, err := range c.uniqueRepoNames(ctx, users, repos) {
		if err != nil {
			return nil, err
		}
		c.Logger.Debug("repo found", "repo", repoName)
		p.Do(func() (github.RepoStats, error) {
			return c.getStats(ctx, repoName)
		})
	}
	return p.Results()
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
