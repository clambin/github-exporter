package github

import (
	"context"
	"github.com/google/go-github/v58/github"
)

type RepoStats struct {
	Name         string
	Stars        int
	Issues       int
	PullRequests int
	Forks        int
	Archived     bool
}

//var _ stats.GitHubClient = Client{}

type Client github.Client

const recordsPerPage = 100

func (c *Client) GetUserReposPage(ctx context.Context, user string, page int) ([]string, int, error) {
	opt := github.RepositoryListByUserOptions{ListOptions: github.ListOptions{Page: page, PerPage: recordsPerPage}}

	repos, resp, err := c.Repositories.ListByUser(ctx, user, &opt)
	if err != nil {
		return nil, 0, err
	}

	repoNames := make([]string, len(repos))
	for i := range repos {
		repoNames[i] = repos[i].GetFullName()
	}

	return repoNames, resp.NextPage, nil
}

func (c *Client) GetRepoStats(ctx context.Context, user string, repo string) (RepoStats, error) {
	var repoStats RepoStats
	r, _, err := c.Repositories.Get(ctx, user, repo)
	if err == nil {
		repoStats.Name = r.GetName()
		repoStats.Stars = r.GetStargazersCount()
		repoStats.Issues = r.GetOpenIssues()
		repoStats.Forks = r.GetForksCount()
		repoStats.Archived = r.GetArchived()
	}
	return repoStats, err
}

func (c *Client) GetPullRequestCountPage(ctx context.Context, user string, repo string, page int) (int, int, error) {
	opt := &github.PullRequestListOptions{ListOptions: github.ListOptions{Page: page, PerPage: recordsPerPage}}
	prs, resp, err := c.PullRequests.List(ctx, user, repo, opt)
	if err != nil {
		return 0, 0, err
	}
	return len(prs), resp.NextPage, nil
}
