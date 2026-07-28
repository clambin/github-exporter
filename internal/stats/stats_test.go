package stats

import (
	"context"
	"log/slog"
	"testing"

	"github.com/clambin/github-exporter/internal/stats/github"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetRepoStats(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		ghClient GitHubClient
		users    []string
		repos    []string
		wantErr  assert.ErrorAssertionFunc
		want     []github.RepoStats
	}{
		{
			name: "success",
			ghClient: fakeGitHubClient{
				userRepoNames: []string{"foo/bar"},
				repoStats:     github.RepoStats{Name: "foo/bar", Stars: 10, Issues: 20, Forks: 1},
				prCount:       5,
			},
			users:   []string{"foo"},
			repos:   nil,
			wantErr: assert.NoError,
			want: []github.RepoStats{
				{Name: "foo/bar", Stars: 10, Issues: 15, PullRequests: 5, Forks: 1},
			},
		},
		{
			name:     "user repos failure",
			ghClient: fakeGitHubClient{err: assert.AnError},
			users:    []string{"foo"},
			repos:    nil,
			wantErr:  assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{GitHubClient: tt.ghClient, Logger: slog.Default()}
			stats, err := c.GetRepoStats(ctx, tt.users, tt.repos)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, stats)
		})
	}
}

func TestClient_getStats(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		ghClient GitHubClient
		repo     string
		wantErr  assert.ErrorAssertionFunc
		want     github.RepoStats
	}{
		{
			name: "success",
			ghClient: fakeGitHubClient{
				repoStats: github.RepoStats{Name: "bar", Stars: 10, Issues: 20, Forks: 1},
				prCount:   5,
			},
			repo:    "foo/bar",
			wantErr: assert.NoError,
			want:    github.RepoStats{Name: "bar", Stars: 10, Issues: 15, PullRequests: 5, Forks: 1},
		},
		{
			name:     "error",
			ghClient: fakeGitHubClient{err: assert.AnError},
			repo:     "foo/bar",
			wantErr:  assert.Error,
		},
		{
			name:     "bad repo name",
			ghClient: fakeGitHubClient{},
			repo:     "foo/bar/snafu",
			wantErr:  assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{GitHubClient: tt.ghClient, Logger: slog.Default()}
			count, err := c.getStats(ctx, tt.repo)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, count)
		})
	}
}

var _ GitHubClient = &fakeGitHubClient{}

type fakeGitHubClient struct {
	userRepoNames []string
	repoStats     github.RepoStats
	prCount       int
	err           error
}

func (f fakeGitHubClient) GetUserRepoNames(_ context.Context, _ string) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.userRepoNames, nil
}

func (f fakeGitHubClient) GetRepoStats(_ context.Context, _ string, _ string) (github.RepoStats, error) {
	if f.err != nil {
		return github.RepoStats{}, f.err
	}
	return f.repoStats, nil
}

func (f fakeGitHubClient) GetPullRequestCount(_ context.Context, _ string, _ string) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.prCount, nil
}
