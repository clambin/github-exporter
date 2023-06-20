package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomnomnom/linkheader"
	"io"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	Token      string
}

const gitHubRepoAPI = "https://api.github.com/"

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return c.HTTPClient.Do(req)
}

func doAllPages[T any](ctx context.Context, c *Client, url string) ([]T, error) {
	var result []T

	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

		resp, err := c.do(req)
		if err != nil {
			return nil, err
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(resp.Status)
		}

		var r []T
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(body, &r)
		}
		if err == nil {
			result = append(result, r...)
		} else {
			return nil, err
		}

		links := resp.Header.Get("Link")
		next := linkheader.Parse(links).FilterByRel("next")
		if len(next) == 0 {
			break
		}

		url = next[0].URL
	}
	return result, nil
}
