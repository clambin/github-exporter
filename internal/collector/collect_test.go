package collector_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/clambin/github-exporter/internal/collector"
	"github.com/clambin/github-exporter/internal/collector/mocks"
	"github.com/clambin/github-exporter/internal/stats/github"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestCollector_Collect(t *testing.T) {
	ctx := context.Background()
	type args struct {
		users           []string
		repos           []string
		includeArchived bool
	}
	tests := []struct {
		name    string
		setup   func(*mocks.StatClient)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    string
	}{
		{
			name: "with archived",
			setup: func(client *mocks.StatClient) {
				client.EXPECT().
					GetRepoStats(ctx, []string{"clambin"}, []string{"clambin/github-exporter", "foo/bar"}).
					Return([]github.RepoStats{
						{Name: "clambin/github-exporter", Stars: 10, Issues: 15, PullRequests: 5, Forks: 1},
						{Name: "clambin/tado-exporter", Stars: 15, Issues: 25, PullRequests: 15, Forks: 2},
						{Name: "foo/bar", Stars: 5, Issues: 5, PullRequests: 5, Forks: 3, Archived: true},
					}, nil)

			},
			args: args{
				users:           []string{"clambin"},
				repos:           []string{"clambin/github-exporter", "foo/bar"},
				includeArchived: true,
			},
			wantErr: assert.NoError,
			want: `
# HELP github_exporter_forks Total number of forks
# TYPE github_exporter_forks gauge
github_exporter_forks{archived="false",repo="clambin/github-exporter"} 1
github_exporter_forks{archived="false",repo="clambin/tado-exporter"} 2
github_exporter_forks{archived="true",repo="foo/bar"} 3
# HELP github_exporter_issues Total number of open issues
# TYPE github_exporter_issues gauge
github_exporter_issues{archived="false",repo="clambin/github-exporter"} 15
github_exporter_issues{archived="false",repo="clambin/tado-exporter"} 25
github_exporter_issues{archived="true",repo="foo/bar"} 5
# HELP github_exporter_pulls Total number of open pull requests
# TYPE github_exporter_pulls gauge
github_exporter_pulls{archived="false",repo="clambin/github-exporter"} 5
github_exporter_pulls{archived="false",repo="clambin/tado-exporter"} 15
github_exporter_pulls{archived="true",repo="foo/bar"} 5
# HELP github_exporter_stars Total number of stars
# TYPE github_exporter_stars gauge
github_exporter_stars{archived="false",repo="clambin/github-exporter"} 10
github_exporter_stars{archived="false",repo="clambin/tado-exporter"} 15
github_exporter_stars{archived="true",repo="foo/bar"} 5
`,
		},
		{
			name: "without archived",
			setup: func(client *mocks.StatClient) {
				client.EXPECT().
					GetRepoStats(ctx, []string{"clambin"}, []string{"clambin/github-exporter", "foo/bar"}).
					Return([]github.RepoStats{
						{Name: "clambin/github-exporter", Stars: 10, Issues: 15, PullRequests: 5, Forks: 1},
						{Name: "clambin/tado-exporter", Stars: 15, Issues: 25, PullRequests: 15, Forks: 2},
						{Name: "foo/bar", Stars: 5, Issues: 5, PullRequests: 5, Forks: 3, Archived: true},
					}, nil)

			},
			args: args{
				users:           []string{"clambin"},
				repos:           []string{"clambin/github-exporter", "foo/bar"},
				includeArchived: false,
			},
			wantErr: assert.NoError,
			want: `
# HELP github_exporter_forks Total number of forks
# TYPE github_exporter_forks gauge
github_exporter_forks{archived="false",repo="clambin/github-exporter"} 1
github_exporter_forks{archived="false",repo="clambin/tado-exporter"} 2
# HELP github_exporter_issues Total number of open issues
# TYPE github_exporter_issues gauge
github_exporter_issues{archived="false",repo="clambin/github-exporter"} 15
github_exporter_issues{archived="false",repo="clambin/tado-exporter"} 25
# HELP github_exporter_pulls Total number of open pull requests
# TYPE github_exporter_pulls gauge
github_exporter_pulls{archived="false",repo="clambin/github-exporter"} 5
github_exporter_pulls{archived="false",repo="clambin/tado-exporter"} 15
# HELP github_exporter_stars Total number of stars
# TYPE github_exporter_stars gauge
github_exporter_stars{archived="false",repo="clambin/github-exporter"} 10
github_exporter_stars{archived="false",repo="clambin/tado-exporter"} 15
`,
		},
		{
			name: "failure",
			setup: func(client *mocks.StatClient) {
				client.EXPECT().
					GetRepoStats(ctx, []string{"clambin"}, []string{"clambin/github-exporter", "foo/bar"}).
					Return(nil, errors.New("fail"))
			},
			args: args{
				users:           []string{"clambin"},
				repos:           []string{"clambin/github-exporter", "foo/bar"},
				includeArchived: false,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gh := mocks.NewStatClient(t)
			tt.setup(gh)
			c := collector.Collector{
				Client:          gh,
				Users:           tt.args.users,
				Repos:           tt.args.repos,
				IncludeArchived: tt.args.includeArchived,
				Lifetime:        time.Second,
				Logger:          slog.Default(),
			}
			tt.wantErr(t, testutil.CollectAndCompare(&c, bytes.NewBufferString(tt.want)))
			tt.wantErr(t, testutil.CollectAndCompare(&c, bytes.NewBufferString(tt.want)))
		})
	}
}
