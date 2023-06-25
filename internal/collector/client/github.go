package client

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"strings"
)

// wrapper for GitHub client

type Client struct {
	*github.Client
}

func (c *Client) GetUserRepos(ctx context.Context, user string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var repos []*github.Repository
	for {
		r, resp, err := c.Client.Repositories.List(ctx, user, opt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, r...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return repos, nil
}

func (c *Client) GetRepo(ctx context.Context, repo string) (*github.Repository, error) {
	user, repo, err := splitFullName(repo)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	r, _, err := c.Client.Repositories.Get(ctx, user, repo)
	return r, err
}

func (c *Client) GetPullRequests(ctx context.Context, repo string) ([]*github.PullRequest, error) {
	var pullRequests []*github.PullRequest
	user, repo, err := splitFullName(repo)
	if err != nil {
		return nil, fmt.Errorf("pull requests: %w", err)
	}

	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		prs, resp, err := c.Client.PullRequests.List(ctx, user, repo, opt)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, prs...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return pullRequests, nil
}

func splitFullName(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo name: %s", repo)
	}
	return parts[0], parts[1], nil
}
