package githubapi

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

// GetRepository returns the repository by given owner and name.
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, errors.Wrap(err, "get repository")
	}
	return repository, nil
}

type GetOrgRepositoriesOptions struct {
	Fork bool
}

// GetOrgRepositories returns the repositories of the given organization.
func (c *Client) GetOrgRepositories(ctx context.Context, organization string, opts GetOrgRepositoriesOptions) ([]*github.Repository, error) {
	typ := "all"
	if opts.Fork {
		typ = "sources"
	}

	repositories, _, err := c.client.Repositories.ListByOrg(ctx, organization, &github.RepositoryListByOrgOptions{
		Type: typ,
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "get organization repositories")
	}
	return repositories, nil
}

type GetUserRepositoriesOptions struct {
	Fork bool
}

// GetUserRepositories returns the repositories of the given user.
func (c *Client) GetUserRepositories(ctx context.Context, user string, opts GetUserRepositoriesOptions) ([]*github.Repository, error) {
	typ := "all"
	if opts.Fork {
		typ = "sources"
	}

	repositories, _, err := c.client.Repositories.List(ctx, user, &github.RepositoryListOptions{
		Type: typ,
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "get user repositories")
	}
	return repositories, nil
}

// GetContributors returns the contributors of the given owner and repository.
// It only returns the first 5 pages.
func (c *Client) GetContributors(ctx context.Context, owner, repo string) ([]*github.Contributor, error) {
	contributorLists := make([]*github.Contributor, 0, 50*5)
	page := 1

	for {
		contributors, _, err := c.client.Repositories.ListContributors(ctx, owner, repo, &github.ListContributorsOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 50,
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "get contributors")
		}

		if len(contributors) == 0 || page == 5 {
			break
		}

		contributorLists = append(contributorLists, contributors...)
		page++
	}

	return contributorLists, nil
}
