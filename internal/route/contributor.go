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

	"github.com/wuhan005/badges/githubapi"
	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/fileutil"
)

func ContributorsBadgeHandler(ctx context.Context) {
	repoNames := ctx.QueryStrings("repos")
	if len(repoNames) == 0 {
		ctx.String(http.StatusNotFound, "empty repositories")
		return
	}

	repos := make(map[string]struct{})
	// Get repositories from query.
	for _, repo := range repoNames {
		if strings.Count(repo, "/") == 1 {
			repos[repo] = struct{}{}
		}
	}

	contributors := make(map[int]*github.Contributor)
	for repo := range repos {
		repoGroup := strings.SplitN(repo, "/", 2)
		repoContributors, err := githubapi.GetContributors(repoGroup[0], repoGroup[1])
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to get contributors: %v", err)
			return
		}

		for _, contributor := range repoContributors {
			contributors[int(contributor.GetID())] = contributor
		}
	}

	// Sort the contributor IDs.
	contributorIDs := make([]int, 0, len(contributors))
	for contributorID := range contributors {
		contributorIDs = append(contributorIDs, contributorID)
	}
	sort.Ints(contributorIDs)

	count := len(contributors)
	width := 1280
	imgSize := 96
	padding := 10

	perLine := width / (imgSize + padding)
	line := int(math.Ceil(float64(count) / float64(perLine)))

	background := imaging.New(width, line*(imgSize+padding), color.NRGBA{})

	index := -1
	for _, contributorID := range contributorIDs {
		index++
		contributor := contributors[contributorID]
		avatarURL := contributor.GetAvatarURL() + "&s=96"

		// Download image to temporary path.
		imagePath := filepath.Join(os.TempDir(), strconv.Itoa(contributorID))
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
