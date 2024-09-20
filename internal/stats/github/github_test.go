package github

import (
	"context"
	"github.com/clambin/github-exporter/internal/stats/github/mocks"
	"github.com/google/go-github/v65/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_GetUserRepoNames(t *testing.T) {
	c := New(http.DefaultTransport)
	r := mocks.NewRepositories(t)
	c.Repositories = r
	ctx := context.Background()

	r.EXPECT().
		ListByUser(
			ctx,
			"user",
			&github.RepositoryListByUserOptions{ListOptions: github.ListOptions{Page: 0, PerPage: recordsPerPage}},
		).
		Return(
			[]*github.Repository{{FullName: constP("user/repo1")}},
			&github.Response{NextPage: 1},
			nil,
		)
	r.EXPECT().
		ListByUser(
			ctx,
			"user",
			&github.RepositoryListByUserOptions{ListOptions: github.ListOptions{Page: 1, PerPage: recordsPerPage}},
		).
		Return(
			[]*github.Repository{{FullName: constP("user/repo2")}},
			&github.Response{NextPage: 0},
			nil,
		)

	repos, err := c.GetUserRepoNames(ctx, "user")
	assert.NoError(t, err)
	assert.Equal(t, []string{"user/repo1", "user/repo2"}, repos)
}

func TestClient_GetRepoStats(t *testing.T) {
	c := New(http.DefaultTransport)

	r := mocks.NewRepositories(t)
	c.Repositories = r

	ctx := context.Background()

	r.EXPECT().Get(ctx, "user", "repo").Return(&github.Repository{
		Owner:           &github.User{Name: constP("user")},
		Name:            constP("repo"),
		ForksCount:      constP(1),
		OpenIssuesCount: constP(2),
		StargazersCount: constP(4),
		Archived:        constP(true),
	}, nil, nil)

	repos, err := c.GetRepoStats(ctx, "user", "repo")
	assert.NoError(t, err)
	assert.Equal(t, RepoStats{
		Name:         "repo",
		Stars:        4,
		Issues:       2,
		PullRequests: 0,
		Forks:        1,
		Archived:     true,
	}, repos)
}

func TestClient_GetPullRequestCount(t *testing.T) {
	c := New(http.DefaultTransport)
	p := mocks.NewPullRequests(t)
	c.PullRequests = p
	ctx := context.Background()

	p.EXPECT().
		List(
			ctx,
			"user",
			"repo",
			&github.PullRequestListOptions{ListOptions: github.ListOptions{Page: 0, PerPage: recordsPerPage}},
		).
		Return(
			[]*github.PullRequest{{}},
			&github.Response{NextPage: 1},
			nil,
		)
	p.EXPECT().
		List(
			ctx,
			"user",
			"repo",
			&github.PullRequestListOptions{ListOptions: github.ListOptions{Page: 1, PerPage: recordsPerPage}},
		).
		Return(
			[]*github.PullRequest{{}},
			&github.Response{NextPage: 0},
			nil,
		)

	prs, err := c.GetPullRequestCount(ctx, "user", "repo")
	assert.NoError(t, err)
	assert.Equal(t, 2, prs)
}

func constP[T any](val T) *T {
	return &val
}
