package stats

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/clambin/github-exporter/internal/stats/github"
	"github.com/clambin/github-exporter/internal/stats/mocks"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetRepoStats(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		setup   func(client *mocks.GitHubClient)
		users   []string
		repos   []string
		wantErr assert.ErrorAssertionFunc
		want    []github.RepoStats
	}{
		{
			name: "success",
			setup: func(client *mocks.GitHubClient) {
				client.EXPECT().GetUserRepoNames(ctx, "foo").Return([]string{"foo/bar"}, nil)
				client.EXPECT().GetRepoStats(ctx, "foo", "bar").Return(github.RepoStats{Name: "foo/bar", Stars: 10, Issues: 20, Forks: 1}, nil)
				client.EXPECT().GetPullRequestCount(ctx, "foo", "bar").Return(5, nil)
			},
			users:   []string{"foo"},
			repos:   nil,
			wantErr: assert.NoError,
			want: []github.RepoStats{
				{Name: "foo/bar", Stars: 10, Issues: 15, PullRequests: 5, Forks: 1},
			},
		},
		{
			name: "user repos failure",
			setup: func(client *mocks.GitHubClient) {
				client.EXPECT().GetUserRepoNames(ctx, "foo").Return(nil, errors.New("fail"))
			},
			users:   []string{"foo"},
			repos:   nil,
			wantErr: assert.Error,
		},
		{
			name: "repo stats failure",
			setup: func(client *mocks.GitHubClient) {
				client.EXPECT().GetUserRepoNames(ctx, "foo").Return([]string{"foo/bar"}, nil)
				client.EXPECT().GetRepoStats(ctx, "foo", "bar").Return(github.RepoStats{Name: "foo/bar", Stars: 10, Issues: 20, Forks: 1}, nil)
				client.EXPECT().GetPullRequestCount(ctx, "foo", "bar").Return(0, errors.New("fail"))
			},
			users:   []string{"foo"},
			repos:   nil,
			wantErr: assert.Error,
			//want:    []github.RepoStats{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gh := mocks.NewGitHubClient(t)
			tt.setup(gh)

			c := Client{GitHubClient: gh, Logger: slog.Default()}

			stats, err := c.GetRepoStats(ctx, tt.users, tt.repos)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, stats)
		})
	}
}

func TestClient_getStats(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		setup   func(*mocks.GitHubClient)
		repo    string
		wantErr assert.ErrorAssertionFunc
		want    github.RepoStats
	}{
		{
			name: "success",
			setup: func(c *mocks.GitHubClient) {
				c.EXPECT().GetRepoStats(ctx, "foo", "bar").Return(github.RepoStats{Name: "bar", Stars: 10, Issues: 20, Forks: 1}, nil)
				c.EXPECT().GetPullRequestCount(ctx, "foo", "bar").Return(5, nil)
			},
			repo:    "foo/bar",
			wantErr: assert.NoError,
			want:    github.RepoStats{Name: "bar", Stars: 10, Issues: 15, PullRequests: 5, Forks: 1},
		},
		{
			name: "error",
			setup: func(c *mocks.GitHubClient) {
				c.EXPECT().GetRepoStats(ctx, "foo", "bar").Return(github.RepoStats{}, errors.New("fail"))
			},
			repo:    "foo/bar",
			wantErr: assert.Error,
		},
		{
			name:    "bad repo name",
			setup:   func(_ *mocks.GitHubClient) {},
			repo:    "foo/bar/snafu",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gh := mocks.NewGitHubClient(t)
			tt.setup(gh)
			c := Client{GitHubClient: gh, Logger: slog.Default()}

			count, err := c.getStats(ctx, tt.repo)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, count)
		})
	}
}
