package collector

import (
	"bytes"
	"errors"
	"github.com/clambin/github-exporter/internal/collector/mocks"
	github2 "github.com/clambin/github-exporter/pkg/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	c := Collector{Client: makeTestClient(t), Users: []string{"clambin"}, Repos: []string{"foo/bar"}}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	buf := bytes.NewBufferString(`# HELP github_monitor_forks Total number of forks
# TYPE github_monitor_forks gauge
github_monitor_forks{archived="false",fork="false",private="false",repo="clambin/github-exporter"} 50
github_monitor_forks{archived="false",fork="false",private="false",repo="clambin/mediamon"} 5
github_monitor_forks{archived="false",fork="false",private="false",repo="foo/bar"} 500
# HELP github_monitor_issues Total number of open issues
# TYPE github_monitor_issues gauge
github_monitor_issues{archived="false",fork="false",private="false",repo="clambin/github-exporter"} 89
github_monitor_issues{archived="false",fork="false",private="false",repo="clambin/mediamon"} 8
github_monitor_issues{archived="false",fork="false",private="false",repo="foo/bar"} 899
# HELP github_monitor_pulls Total number of open pull requests
# TYPE github_monitor_pulls gauge
github_monitor_pulls{archived="false",fork="false",private="false",repo="clambin/github-exporter"} 1
github_monitor_pulls{archived="false",fork="false",private="false",repo="clambin/mediamon"} 1
github_monitor_pulls{archived="false",fork="false",private="false",repo="foo/bar"} 1
# HELP github_monitor_stars Total number of stars
# TYPE github_monitor_stars gauge
github_monitor_stars{archived="false",fork="false",private="false",repo="clambin/github-exporter"} 100
github_monitor_stars{archived="false",fork="false",private="false",repo="clambin/mediamon"} 10
github_monitor_stars{archived="false",fork="false",private="false",repo="foo/bar"} 1000
`)
	assert.NoError(t, testutil.GatherAndCompare(r, buf))
}

func TestCollector_Collect_Failure(t *testing.T) {
	client := mocks.NewGitHubClient(t)
	client.On("GetUserRepos", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, errors.New("failure"))

	c := Collector{Client: client, Users: []string{"clambin"}}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	_, err := r.Gather()
	assert.Error(t, err)
}

func makeTestClient(t *testing.T) GitHubClient {
	client := mocks.NewGitHubClient(t)
	client.
		On("GetUserRepos", mock.AnythingOfType("*context.emptyCtx"), "clambin").
		Return([]github2.Repo{
			{
				FullName:        "clambin/mediamon",
				StargazersCount: 10,
				ForksCount:      5,
				OpenIssuesCount: 9,
			},
			{
				FullName:        "clambin/github-exporter",
				StargazersCount: 100,
				ForksCount:      50,
				OpenIssuesCount: 90,
			},
			{
				FullName:        "clambin/snafu",
				Archived:        true,
				StargazersCount: 100,
				ForksCount:      50,
				OpenIssuesCount: 90,
			},
		}, nil)
	client.
		On("GetRepo", mock.AnythingOfType("*context.emptyCtx"), "foo/bar").
		Return(github2.Repo{
			FullName:        "foo/bar",
			StargazersCount: 1000,
			ForksCount:      500,
			OpenIssuesCount: 900,
		}, nil)
	client.
		On("GetPullRequests", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Return([]github2.PullRequest{{}}, nil)

	return client
}
