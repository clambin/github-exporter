package github

import (
	"context"
	"time"
)

func (c *Client) GetIssues(ctx context.Context, repo string) ([]Issue, error) {
	url := gitHubRepoAPI + "repos/" + repo + "/issues"
	return doAllPages[Issue](ctx, c, url)
}

type Issue struct {
	Url           string `json:"url"`
	RepositoryUrl string `json:"repository_url"`
	LabelsUrl     string `json:"labels_url"`
	CommentsUrl   string `json:"comments_url"`
	EventsUrl     string `json:"events_url"`
	HtmlUrl       string `json:"html_url"`
	Id            int    `json:"id"`
	NodeId        string `json:"node_id"`
	Number        int    `json:"number"`
	Title         string `json:"title"`
	User          struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"user"`
	Labels            []interface{} `json:"labels"`
	State             string        `json:"state"`
	Locked            bool          `json:"locked"`
	Assignee          interface{}   `json:"assignee"`
	Assignees         []interface{} `json:"assignees"`
	Milestone         interface{}   `json:"milestone"`
	Comments          int           `json:"comments"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	ClosedAt          interface{}   `json:"closed_at"`
	AuthorAssociation string        `json:"author_association"`
	ActiveLockReason  interface{}   `json:"active_lock_reason"`
	Draft             bool          `json:"draft,omitempty"`
	PullRequest       struct {
		Url      string      `json:"url"`
		HtmlUrl  string      `json:"html_url"`
		DiffUrl  string      `json:"diff_url"`
		PatchUrl string      `json:"patch_url"`
		MergedAt interface{} `json:"merged_at"`
	} `json:"pull_request,omitempty"`
	Body      *string `json:"body"`
	Reactions struct {
		Url        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Field3     int    `json:"+1"`
		Field4     int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions"`
	TimelineUrl           string      `json:"timeline_url"`
	PerformedViaGithubApp interface{} `json:"performed_via_github_app"`
	StateReason           interface{} `json:"state_reason"`
}
