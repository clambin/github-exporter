package github

import (
	"context"
	"encoding/json"
	"errors"
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

func doOnePage[T any](ctx context.Context, c *Client, url string) ([]T, string, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := c.do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New(resp.Status)
	}

	var next string
	links := resp.Header.Get("Link")
	if nextLinks := linkheader.Parse(links).FilterByRel("next"); len(nextLinks) > 0 {
		next = nextLinks[0].URL
	}

	var r []T
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &r)
	}

	return r, next, err
}

func doAllPages[T any](ctx context.Context, c *Client, url string) ([]T, error) {
	var result, page []T
	var err error
	for {
		if page, url, err = doOnePage[T](ctx, c, url); err == nil {
			result = append(result, page...)
		}

		if err != nil || url == "" {
			break
		}
	}
	return result, err
}
