package collector_test

import (
	"bytes"
	"errors"
	"github.com/clambin/github-exporter/internal/collector"
	"github.com/clambin/github-exporter/internal/collector/mocks"
	"github.com/google/go-github/v53/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCollector_Collect(t *testing.T) {
	c := collector.Collector{
		collector.GitHubCache{
			Client:          makeTestClient(t),
			Users:           []string{"clambin"},
			Repos:           []string{"foo/bar"},
			IncludeArchived: false,
			Lifetime:        time.Hour,
		},
	}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	expected := `# HELP github_monitor_forks Total number of forks
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
`
	assert.NoError(t, testutil.GatherAndCompare(r, bytes.NewBufferString(expected)))
	assert.NoError(t, testutil.GatherAndCompare(r, bytes.NewBufferString(expected)))
}

func TestCollector_Collect_Failure(t *testing.T) {
	client := mocks.NewGitHubClient(t)
	client.On("GetUserRepos", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, errors.New("failure"))

	c := collector.Collector{
		collector.GitHubCache{
			Client: client,
			Users:  []string{"clambin"},
		},
	}
	r := prometheus.NewPedanticRegistry()
	r.MustRegister(&c)

	_, err := r.Gather()
	assert.Error(t, err)
}

func makeTestClient(t *testing.T) collector.GitHubClient {
	repos := []struct {
		name     string
		archived bool
		stars    int
		forks    int
		issues   int
	}{
		{name: "clambin/mediamon", stars: 10, forks: 5, issues: 9},
		{name: "clambin/github-exporter", stars: 100, forks: 50, issues: 90},
		{name: "clambin/snafu", archived: true, stars: 100, forks: 50, issues: 90},
		{name: "foo/bar", stars: 1000, forks: 500, issues: 900},
	}

	client := mocks.NewGitHubClient(t)
	client.
		On("GetUserRepos", mock.AnythingOfType("*context.emptyCtx"), "clambin").
		Return([]*github.Repository{
			{
				FullName:        &repos[0].name,
				Archived:        &repos[0].archived,
				StargazersCount: &repos[0].stars,
				ForksCount:      &repos[0].forks,
				OpenIssuesCount: &repos[0].issues,
			},
			{
				FullName:        &repos[1].name,
				Archived:        &repos[1].archived,
				StargazersCount: &repos[1].stars,
				ForksCount:      &repos[1].forks,
				OpenIssuesCount: &repos[1].issues,
			},
			{
				FullName:        &repos[2].name,
				Archived:        &repos[2].archived,
				StargazersCount: &repos[2].stars,
				ForksCount:      &repos[2].forks,
				OpenIssuesCount: &repos[2].issues,
			},
		}, nil).Once()
	client.
		On("GetRepo", mock.AnythingOfType("*context.emptyCtx"), "foo/bar").
		Return(&github.Repository{
			FullName:        &repos[3].name,
			Archived:        &repos[3].archived,
			StargazersCount: &repos[3].stars,
			ForksCount:      &repos[3].forks,
			OpenIssuesCount: &repos[3].issues,
		}, nil).Once()
	client.
		On("GetPullRequests", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Return([]*github.PullRequest{{}}, nil)

	return client
}
