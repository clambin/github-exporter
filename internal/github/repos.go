package github

import (
	"context"
)

func (c *Client) GetUserRepos(ctx context.Context, owner string) ([]Repo, error) {
	url := gitHubRepoAPI + "users/" + owner + "/repos"
	return doAllPages[Repo](ctx, c, url)
}
