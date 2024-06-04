package github

import (
	"context"
	"github.com/google/go-github/v62/github"
	"net/http"
	"time"
)

type RepoStats struct {
	Name         string
	Stars        int
	Issues       int
	PullRequests int
	Forks        int
	Archived     bool
}

type Client struct {
	Repositories
	PullRequests
}

type Repositories interface {
	ListByUser(context.Context, string, *github.RepositoryListByUserOptions) ([]*github.Repository, *github.Response, error)
	Get(context.Context, string, string) (*github.Repository, *github.Response, error)
}

type PullRequests interface {
	List(context.Context, string, string, *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

func New(tp http.RoundTripper) *Client {
	client := github.NewClient(&http.Client{Transport: tp, Timeout: 10 * time.Second})
	return &Client{
		Repositories: client.Repositories,
		PullRequests: client.PullRequests,
	}
}

const recordsPerPage = 100

func (c Client) GetUserRepoNames(ctx context.Context, user string) (repos []string, err error) {
	var page int
	for err == nil {
		var repoPage []string
		if repoPage, page, err = c.GetUserReposPage(ctx, user, page); err == nil {
			repos = append(repos, repoPage...)
			if page == 0 {
				break
			}
		}
	}
	return repos, err
}

func (c Client) GetUserReposPage(ctx context.Context, user string, page int) (repoNames []string, nextPage int, err error) {
	opt := github.RepositoryListByUserOptions{ListOptions: github.ListOptions{Page: page, PerPage: recordsPerPage}}

	var repos []*github.Repository
	var resp *github.Response
	if repos, resp, err = c.Repositories.ListByUser(ctx, user, &opt); err == nil {
		repoNames = make([]string, len(repos))
		for i := range repos {
			repoNames[i] = repos[i].GetFullName()
		}
		nextPage = resp.NextPage
	}

	return repoNames, nextPage, err
}

func (c Client) GetRepoStats(ctx context.Context, user string, repo string) (RepoStats, error) {
	var repoStats RepoStats
	r, _, err := c.Repositories.Get(ctx, user, repo)
	if err == nil {
		repoStats.Name = r.GetName()
		repoStats.Stars = r.GetStargazersCount()
		repoStats.Issues = r.GetOpenIssuesCount()
		repoStats.Forks = r.GetForksCount()
		repoStats.Archived = r.GetArchived()
	}
	return repoStats, err
}

func (c Client) GetPullRequestCount(ctx context.Context, user string, repo string) (prCount int, err error) {
	var page int
	for err == nil {
		var count int
		if count, page, err = c.GetPullRequestCountPage(ctx, user, repo, page); err == nil {
			prCount += count
			if page == 0 {
				break
			}
		}
	}
	return prCount, err
}

func (c Client) GetPullRequestCountPage(ctx context.Context, user string, repo string, page int) (prCount int, nextPage int, err error) {
	opt := &github.PullRequestListOptions{ListOptions: github.ListOptions{Page: page, PerPage: recordsPerPage}}
	var prs []*github.PullRequest
	var resp *github.Response
	if prs, resp, err = c.PullRequests.List(ctx, user, repo, opt); err == nil {
		prCount = len(prs)
		nextPage = resp.NextPage
	}
	return prCount, nextPage, err
}
