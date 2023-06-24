package github

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestClient_GetRepos(t *testing.T) {
	c := Client{
		HTTPClient: &http.Client{
			Transport: &transport{path.Join("testdata", "repos.json")},
		},
	}
	ctx := context.Background()
	repos, err := c.GetUserRepos(ctx, "clambin")
	require.NoError(t, err)
	assert.Len(t, repos, 30)
}

func TestClient_GetRepo(t *testing.T) {
	c := Client{
		HTTPClient: &http.Client{
			Transport: &transport{path.Join("testdata", "repo.json")},
		},
	}
	ctx := context.Background()
	repo, err := c.GetRepo(ctx, "clambin/tado-exporter")
	require.NoError(t, err)

	assert.Equal(t, "tado-exporter", repo.Name)
	assert.Equal(t, 6, repo.StargazersCount)
	assert.Equal(t, "public", repo.Visibility)
}

func TestClient_GetPullRequests(t *testing.T) {
	c := Client{
		HTTPClient: &http.Client{
			Transport: &transport{path.Join("testdata", "pulls.json")},
		},
	}
	ctx := context.Background()
	pulls, err := c.GetPullRequests(ctx, "githubexporter/github-exporter")
	require.NoError(t, err)
	assert.Len(t, pulls, 10)
}

type transport struct {
	responseFile string
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, mandatoryHeader := range []string{"Accept", "Authorization", "X-GitHub-Api-Version"} {
		if req.Header.Get(mandatoryHeader) == "" {
			return nil, fmt.Errorf("missing mandatory header: %s", mandatoryHeader)
		}
	}
	response, err := os.ReadFile(t.responseFile)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		Status:     "OK",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(response)),
	}, nil
}
