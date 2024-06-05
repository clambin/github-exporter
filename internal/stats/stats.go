package stats

import (
	"context"
	"fmt"
	"github.com/clambin/github-exporter/internal/stats/github"
	"github.com/clambin/go-common/set"
	"log/slog"
	"strings"
	"time"
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
	uniqueRepos, err := c.getUniqueRepoNames(ctx, users, repos)
	if err != nil {
		return nil, err
	}
	c.Logger.Debug("unique repos", "repos", uniqueRepos)

	var p parallel[github.RepoStats]

	for i := range uniqueRepos {
		p.Do(func() (github.RepoStats, error) {
			return c.getStats(ctx, uniqueRepos[i])
		})
	}
	return p.Results()
}

func (c Client) getUniqueRepoNames(ctx context.Context, users []string, repos []string) ([]string, error) {
	uniqueRepoNames := set.New(repos...)

	for _, user := range users {
		userRepos, err := c.GitHubClient.GetUserRepoNames(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("get repos for user %s: %w", user, err)
		}
		uniqueRepoNames.Add(userRepos...)
	}

	return uniqueRepoNames.List(), nil
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
	repoStats.PullRequests, err = c.GitHubClient.GetPullRequestCount(ctx, user, repo)
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
