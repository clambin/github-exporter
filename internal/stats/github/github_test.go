package github

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v89/github"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetUserRepoNames(t *testing.T) {
	c, _ := New(http.DefaultTransport)
	c.Repositories = fakeRepositories{
		repoList: map[int]repoPage{
			0: {
				repo: []*github.Repository{{FullName: new("user/repo1")}},
				resp: &github.Response{NextPage: 1},
			},
			1: {
				repo: []*github.Repository{{FullName: new("user/repo2")}},
				resp: &github.Response{NextPage: 0},
			},
		},
	}
	ctx := context.Background()

	repos, err := c.GetUserRepoNames(ctx, "user")
	assert.NoError(t, err)
	assert.Equal(t, []string{"user/repo1", "user/repo2"}, repos)
}

func TestClient_GetRepoStats(t *testing.T) {
	c, _ := New(http.DefaultTransport)
	c.Repositories = fakeRepositories{
		repos: map[string]*github.Repository{
			"user/repo": {
				Owner:           &github.User{Name: new("user")},
				Name:            new("repo"),
				ForksCount:      new(1),
				OpenIssuesCount: new(2),
				StargazersCount: new(4),
				Archived:        new(true),
			},
		},
	}

	ctx := context.Background()
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
	c, _ := New(http.DefaultTransport)
	p := fakePullRequests{
		prs: map[int]prPage{
			0: {
				prs:  []*github.PullRequest{{}},
				resp: &github.Response{NextPage: 1},
			},
			1: {
				prs:  []*github.PullRequest{{}},
				resp: &github.Response{NextPage: 0},
			},
		},
	}
	c.PullRequests = p
	ctx := context.Background()

	prs, err := c.GetPullRequestCount(ctx, "user", "repo")
	assert.NoError(t, err)
	assert.Equal(t, 2, prs)
}

var _ Repositories = &fakeRepositories{}

type repoPage struct {
	repo []*github.Repository
	resp *github.Response
}
type fakeRepositories struct {
	repoList map[int]repoPage
	repos    map[string]*github.Repository
}

func (f fakeRepositories) ListByUser(_ context.Context, _ string, options *github.RepositoryListByUserOptions) ([]*github.Repository, *github.Response, error) {
	page, found := f.repoList[options.Page]
	if !found {
		return nil, nil, errors.New("page not found")
	}
	return page.repo, page.resp, nil
}

func (f fakeRepositories) Get(ctx context.Context, s1 string, s2 string) (*github.Repository, *github.Response, error) {
	repo, found := f.repos[s1+"/"+s2]
	if !found {
		return nil, nil, errors.New("repo not found")
	}
	return repo, &github.Response{}, nil
}

var _ PullRequests = &fakePullRequests{}

type prPage struct {
	prs  []*github.PullRequest
	resp *github.Response
}
type fakePullRequests struct {
	prs map[int]prPage
}

func (f fakePullRequests) List(_ context.Context, _ string, _ string, options *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	page, found := f.prs[options.Page]
	if !found {
		return nil, nil, errors.New("prs not found")
	}
	return page.prs, page.resp, nil
}
