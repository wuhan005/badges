package githubapi

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

func GetContributors(owner, repo string) ([]*github.Contributor, error) {
	client := github.NewClient(nil)
	ctx := context.Background()

	contributorLists := make([]*github.Contributor, 0, 50*5)
	page := 1

	for {
		contributors, _, err := client.Repositories.ListContributors(ctx, owner, repo, &github.ListContributorsOptions{
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
