package route

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/go-github/v29/github"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/badges/githubapi"
	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/fileutil"
)

func ContributorsBadgeHandler(ctx context.Context) {
	orgNames := ctx.QueryStrings("orgs")
	repoFullNames := ctx.QueryStrings("repos")
	fork := ctx.QueryBool("fork")
	if len(orgNames) == 0 && len(repoFullNames) == 0 {
		ctx.String(http.StatusNotFound, "both organization and repository is empty")
		return
	}

	client := githubapi.NewGitHubAPIClient()
	reposSet := make(map[string]*github.Repository)
	repos := make([]*github.Repository, 0)

	// Get the repositories of the organizations.
	orgsSet := make(map[string]struct{})
	for _, orgName := range orgNames {
		if _, ok := orgsSet[orgName]; ok {
			continue
		}

		orgRepos, err := client.GetOrgRepositories(ctx.Request().Context(), orgName,
			githubapi.GetOrgRepositoriesOptions{
				Fork: fork,
			},
		)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to get organization's repositories: %v", err)
			return
		}

		log.Trace("Get %d repositories of org %s", len(orgRepos), orgName)

		for _, orgRepo := range orgRepos {
			repo := orgRepo
			if _, ok := reposSet[*repo.FullName]; ok {
				continue
			}
			reposSet[*repo.FullName] = repo
			repos = append(repos, repo)
		}
	}

	// Get the repositories specified by user.
	for _, repoFullName := range repoFullNames {
		if _, ok := reposSet[repoFullName]; ok {
			continue
		}

		if strings.Count(repoFullName, "/") == 1 {
			repoGroup := strings.SplitN(repoFullName, "/", 2)
			owner, repoName := repoGroup[0], repoGroup[1]

			repo, err := client.GetRepository(ctx.Request().Context(), owner, repoName)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to get repository: %v", err)
				return
			}

			reposSet[*repo.FullName] = repo
			repos = append(repos, repo)
		}
	}

	// Sort repositories by stargazers count.
	sort.SliceStable(repos, func(i, j int) bool {
		// We need the descending order.
		return *repos[i].StargazersCount > *repos[j].StargazersCount
	})

	// Get contributors.
	contributorsSet := make(map[int64]struct{})
	contributors := make([]*github.Contributor, 0)
	for _, repo := range repos {
		repoContributors, err := client.GetContributors(ctx.Request().Context(), *repo.Owner.Login, *repo.Name)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to get contributors: %v", err)
			return
		}

		log.Trace("Get %d contributors of repo %s (%d stars)", len(repoContributors), *repo.FullName, *repo.StargazersCount)

		for _, contributor := range repoContributors {
			if _, ok := contributorsSet[*contributor.ID]; ok {
				continue
			}
			contributors = append(contributors, contributor)
			contributorsSet[*contributor.ID] = struct{}{}
		}
	}

	count := len(contributors)
	width := 1280
	imgSize := 96
	padding := 10

	perLine := width / (imgSize + padding)
	line := int(math.Ceil(float64(count) / float64(perLine)))

	background := imaging.New(width, line*(imgSize+padding), color.NRGBA{})

	index := -1
	for _, contributor := range contributors {
		index++
		avatarURL := contributor.GetAvatarURL() + "&s=96"

		// Download image to temporary path.
		imagePath := filepath.Join(os.TempDir(), strconv.Itoa(int(*contributor.ID)))
		err := fileutil.DownloadFile(avatarURL, imagePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to download contributor's avatar: %v", err)
			return
		}

		// Load image file
		img, err := imaging.Open(imagePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to load image file: %v", err)
			return
		}

		// Resize image
		img = imaging.Resize(img, imgSize, imgSize, imaging.Lanczos)

		// Paste image into background
		y := int(math.Floor(float64((index)/perLine))) * (imgSize + padding)
		x := (index - int(math.Floor(float64((index)/perLine)))*perLine) * (imgSize + padding)
		background = imaging.Paste(background, img, image.Pt(x, y))
	}

	encoder := png.Encoder{}
	err := encoder.Encode(ctx.ResponseWriter(), background)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to encode image: %v", err)
		return
	}
}
