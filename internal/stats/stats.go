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
	GetUserReposPage(context.Context, string, int) ([]string, int, error)
	GetRepoStats(context.Context, string, string) (github.RepoStats, error)
	GetPullRequestCountPage(context.Context, string, string, int) (int, int, error)
}

func (c Client) GetRepoStats(ctx context.Context, users []string, repos []string) ([]github.RepoStats, error) {
	uniqueRepos, err := c.getUniqueRepoNames(ctx, users, repos)
	if err != nil {
		return nil, err
	}
	c.Logger.Debug("unique repos", "repos", uniqueRepos)

	var p parallel[github.RepoStats]

	for i := range uniqueRepos {
		i := i
		p.Do(func() (github.RepoStats, error) {
			return c.getStats(ctx, uniqueRepos[i])
		})
	}
	return p.Results()
}

func (c Client) getUniqueRepoNames(ctx context.Context, users []string, repos []string) ([]string, error) {
	unique := set.New(repos...)

	for _, user := range users {
		userRepos, err := c.getUserRepoNames(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("get repos for user %s: %w", user, err)
		}
		unique.Add(userRepos...)
	}

	return unique.List(), nil
}

func (c Client) getUserRepoNames(ctx context.Context, user string) ([]string, error) {
	var repos []string
	var page int
	for {
		repoPage, nextPage, err := c.GitHubClient.GetUserReposPage(ctx, user, page)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repoPage...)
		if nextPage == 0 {
			break
		}
		page = nextPage
	}
	return repos, nil
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
	repoStats.PullRequests, err = c.getPullRequestCount(ctx, user, repo)
	if err != nil {
		return repoStats, err
	}
	repoStats.Issues -= repoStats.PullRequests
	return repoStats, nil
}

func (c Client) getPullRequestCount(ctx context.Context, user, repo string) (int, error) {
	var pullRequestCount int
	var page int
	for {
		count, nextPage, err := c.GitHubClient.GetPullRequestCountPage(ctx, user, repo, page)
		if err != nil {
			return 0, err
		}
		pullRequestCount += count
		if nextPage == 0 {
			break
		}
		page = nextPage
	}
	return pullRequestCount, nil
}

func splitFullName(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo name: %s", repo)
	}
	return parts[0], parts[1], nil
}
