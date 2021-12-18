// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package githubapi

import (
	"net/http"

	"github.com/google/go-github/v29/github"
)

type Client struct {
	client *github.Client
}

// NewGitHubAPIClient creates and returns a new GitHub API client.
func NewGitHubAPIClient(httpClients ...*http.Client) *Client {
	httpClient := (*http.Client)(nil)
	if len(httpClients) == 1 {
		httpClient = httpClients[0]
	}

	return &Client{
		client: github.NewClient(httpClient),
	}
}
